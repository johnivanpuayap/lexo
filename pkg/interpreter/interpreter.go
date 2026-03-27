package interpreter

import (
	"fmt"
	"math"
	"strings"

	"github.com/johnivanpuayap/lexo/pkg/checker"
	"github.com/johnivanpuayap/lexo/pkg/parser"
)

type IOHandler interface {
	Print(text string)
	Input(prompt string) string
}

type RuntimeError struct {
	Message string
	Line    int
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("Line %d: %s", e.Line, e.Message)
}

// Sentinel errors for control flow
type breakSignal struct{}
type continueSignal struct{}
type returnSignal struct {
	Value Value
}

func (b breakSignal) Error() string    { return "break" }
func (c continueSignal) Error() string { return "continue" }
func (r returnSignal) Error() string   { return "return" }

type Interpreter struct {
	env       *Environment
	globalEnv *Environment
	io        IOHandler
	funcs     map[string]*parser.FuncDecl
}

func New(io IOHandler) *Interpreter {
	env := NewEnvironment(nil)
	return &Interpreter{
		env:       env,
		globalEnv: env,
		io:        io,
		funcs:     make(map[string]*parser.FuncDecl),
	}
}

func (interp *Interpreter) Execute(prog *parser.Program) error {
	// First pass: collect function declarations
	for _, stmt := range prog.Body {
		if fd, ok := stmt.(*parser.FuncDecl); ok {
			interp.funcs[fd.Name] = fd
		}
	}

	// Second pass: execute top-level statements
	for _, stmt := range prog.Body {
		if _, ok := stmt.(*parser.FuncDecl); ok {
			continue // skip function declarations
		}
		if err := interp.execStmt(stmt); err != nil {
			if _, ok := err.(*RuntimeError); ok {
				return err
			}
			return err
		}
	}
	return nil
}

func (interp *Interpreter) execStmt(stmt parser.Statement) error {
	switch s := stmt.(type) {
	case *parser.VarDecl:
		val, err := interp.evalExpr(s.Value)
		if err != nil {
			return err
		}
		interp.env.Define(s.Name, val)
		return nil

	case *parser.Assignment:
		val, err := interp.evalExpr(s.Value)
		if err != nil {
			return err
		}
		switch target := s.Target.(type) {
		case *parser.Identifier:
			if !interp.env.Assign(target.Name, val) {
				return &RuntimeError{Message: fmt.Sprintf("variable '%s' not defined", target.Name), Line: s.Line}
			}
		case *parser.IndexAccess:
			return interp.assignIndex(target, val, s.Line)
		}
		return nil

	case *parser.ExpressionStmt:
		_, err := interp.evalExpr(s.Expr)
		return err

	case *parser.IfStmt:
		return interp.execIf(s)

	case *parser.WhileStmt:
		return interp.execWhile(s)

	case *parser.ForStmt:
		return interp.execFor(s)

	case *parser.ReturnStmt:
		if s.Value == nil {
			return returnSignal{Value: VoidVal{}}
		}
		val, err := interp.evalExpr(s.Value)
		if err != nil {
			return err
		}
		return returnSignal{Value: val}

	case *parser.BreakStmt:
		return breakSignal{}

	case *parser.ContinueStmt:
		return continueSignal{}

	case *parser.FuncDecl:
		return nil // already collected

	default:
		return &RuntimeError{Message: fmt.Sprintf("unknown statement type %T", stmt), Line: stmt.GetLine()}
	}
}

func (interp *Interpreter) assignIndex(target *parser.IndexAccess, val Value, line int) error {
	obj, err := interp.evalExpr(target.Object)
	if err != nil {
		return err
	}
	arr, ok := obj.(*ArrayVal)
	if !ok {
		return &RuntimeError{Message: "cannot index non-array", Line: line}
	}
	idx, err := interp.evalExpr(target.Index)
	if err != nil {
		return err
	}
	i := int(idx.(IntVal))
	if i < 0 || i >= len(arr.Elements) {
		return &RuntimeError{Message: fmt.Sprintf("index out of bounds: %d (array length %d)", i, len(arr.Elements)), Line: line}
	}
	arr.Elements[i] = val
	return nil
}

func (interp *Interpreter) execIf(s *parser.IfStmt) error {
	cond, err := interp.evalExpr(s.Condition)
	if err != nil {
		return err
	}
	if bool(cond.(BoolVal)) {
		return interp.execBlock(s.Body)
	}

	for _, elif := range s.Elifs {
		cond, err := interp.evalExpr(elif.Condition)
		if err != nil {
			return err
		}
		if bool(cond.(BoolVal)) {
			return interp.execBlock(elif.Body)
		}
	}

	if s.ElseBody != nil {
		return interp.execBlock(s.ElseBody)
	}
	return nil
}

func (interp *Interpreter) execWhile(s *parser.WhileStmt) error {
	for {
		cond, err := interp.evalExpr(s.Condition)
		if err != nil {
			return err
		}
		if !bool(cond.(BoolVal)) {
			break
		}
		err = interp.execBlock(s.Body)
		if err != nil {
			if _, ok := err.(breakSignal); ok {
				break
			}
			if _, ok := err.(continueSignal); ok {
				continue
			}
			return err
		}
	}
	return nil
}

func (interp *Interpreter) execFor(s *parser.ForStmt) error {
	iterable, err := interp.evalExpr(s.Iterable)
	if err != nil {
		return err
	}
	arr := iterable.(*ArrayVal)

	for _, elem := range arr.Elements {
		prevEnv := interp.env
		interp.env = NewEnvironment(interp.env)
		interp.env.Define(s.VarName, elem)

		err := interp.execBlock(s.Body)
		interp.env = prevEnv

		if err != nil {
			if _, ok := err.(breakSignal); ok {
				break
			}
			if _, ok := err.(continueSignal); ok {
				continue
			}
			return err
		}
	}
	return nil
}

func (interp *Interpreter) execBlock(stmts []parser.Statement) error {
	for _, stmt := range stmts {
		if err := interp.execStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (interp *Interpreter) evalExpr(expr parser.Expression) (Value, error) {
	switch e := expr.(type) {
	case *parser.Literal:
		return interp.evalLiteral(e), nil
	case *parser.Identifier:
		val, ok := interp.env.Get(e.Name)
		if !ok {
			return nil, &RuntimeError{Message: fmt.Sprintf("variable '%s' not defined", e.Name), Line: e.Line}
		}
		return val, nil
	case *parser.BinaryExpr:
		return interp.evalBinary(e)
	case *parser.UnaryExpr:
		return interp.evalUnary(e)
	case *parser.FuncCall:
		return interp.evalFuncCall(e)
	case *parser.ArrayLiteral:
		return interp.evalArrayLiteral(e)
	case *parser.IndexAccess:
		return interp.evalIndexAccess(e)
	case *parser.MethodCall:
		return interp.evalMethodCall(e)
	default:
		return nil, &RuntimeError{Message: fmt.Sprintf("unknown expression %T", expr), Line: expr.GetLine()}
	}
}

func (interp *Interpreter) evalLiteral(e *parser.Literal) Value {
	switch e.LitType {
	case "int":
		return IntVal(e.Value.(int64))
	case "float":
		return FloatVal(e.Value.(float64))
	case "string":
		return StringVal(e.Value.(string))
	case "bool":
		return BoolVal(e.Value.(bool))
	default:
		return VoidVal{}
	}
}

func (interp *Interpreter) evalBinary(e *parser.BinaryExpr) (Value, error) {
	// Short-circuit for logical operators
	if e.Operator == "and" || e.Operator == "or" {
		left, err := interp.evalExpr(e.Left)
		if err != nil {
			return nil, err
		}
		leftBool := bool(left.(BoolVal))
		if e.Operator == "and" && !leftBool {
			return BoolVal(false), nil
		}
		if e.Operator == "or" && leftBool {
			return BoolVal(true), nil
		}
		right, err := interp.evalExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return BoolVal(bool(right.(BoolVal))), nil
	}

	left, err := interp.evalExpr(e.Left)
	if err != nil {
		return nil, err
	}
	right, err := interp.evalExpr(e.Right)
	if err != nil {
		return nil, err
	}

	switch e.Operator {
	case "+":
		switch l := left.(type) {
		case IntVal:
			return IntVal(int64(l) + int64(right.(IntVal))), nil
		case FloatVal:
			return FloatVal(float64(l) + float64(right.(FloatVal))), nil
		case StringVal:
			return StringVal(string(l) + string(right.(StringVal))), nil
		}
	case "-":
		switch l := left.(type) {
		case IntVal:
			return IntVal(int64(l) - int64(right.(IntVal))), nil
		case FloatVal:
			return FloatVal(float64(l) - float64(right.(FloatVal))), nil
		}
	case "*":
		switch l := left.(type) {
		case IntVal:
			return IntVal(int64(l) * int64(right.(IntVal))), nil
		case FloatVal:
			return FloatVal(float64(l) * float64(right.(FloatVal))), nil
		}
	case "/":
		switch l := left.(type) {
		case IntVal:
			r := int64(right.(IntVal))
			if r == 0 {
				return nil, &RuntimeError{Message: "division by zero", Line: e.Line}
			}
			return IntVal(int64(l) / r), nil
		case FloatVal:
			r := float64(right.(FloatVal))
			if r == 0 {
				return nil, &RuntimeError{Message: "division by zero", Line: e.Line}
			}
			return FloatVal(float64(l) / r), nil
		}
	case "%":
		switch l := left.(type) {
		case IntVal:
			r := int64(right.(IntVal))
			if r == 0 {
				return nil, &RuntimeError{Message: "division by zero", Line: e.Line}
			}
			return IntVal(int64(l) % r), nil
		case FloatVal:
			return FloatVal(math.Mod(float64(l), float64(right.(FloatVal)))), nil
		}
	case "==":
		return BoolVal(interp.equals(left, right)), nil
	case "!=":
		return BoolVal(!interp.equals(left, right)), nil
	case "<":
		return interp.compare(left, right, "<", e.Line)
	case ">":
		return interp.compare(left, right, ">", e.Line)
	case "<=":
		return interp.compare(left, right, "<=", e.Line)
	case ">=":
		return interp.compare(left, right, ">=", e.Line)
	// "and" and "or" are handled above with short-circuit evaluation
	}

	return nil, &RuntimeError{Message: fmt.Sprintf("unsupported binary operation: %s", e.Operator), Line: e.Line}
}

func (interp *Interpreter) equals(a, b Value) bool {
	switch av := a.(type) {
	case IntVal:
		return int64(av) == int64(b.(IntVal))
	case FloatVal:
		return float64(av) == float64(b.(FloatVal))
	case StringVal:
		return string(av) == string(b.(StringVal))
	case BoolVal:
		return bool(av) == bool(b.(BoolVal))
	}
	return false
}

func (interp *Interpreter) compare(a, b Value, op string, line int) (Value, error) {
	switch av := a.(type) {
	case IntVal:
		bv := int64(b.(IntVal))
		ai := int64(av)
		switch op {
		case "<":
			return BoolVal(ai < bv), nil
		case ">":
			return BoolVal(ai > bv), nil
		case "<=":
			return BoolVal(ai <= bv), nil
		case ">=":
			return BoolVal(ai >= bv), nil
		}
	case FloatVal:
		bv := float64(b.(FloatVal))
		af := float64(av)
		switch op {
		case "<":
			return BoolVal(af < bv), nil
		case ">":
			return BoolVal(af > bv), nil
		case "<=":
			return BoolVal(af <= bv), nil
		case ">=":
			return BoolVal(af >= bv), nil
		}
	}
	return nil, &RuntimeError{Message: "cannot compare values", Line: line}
}

func (interp *Interpreter) evalUnary(e *parser.UnaryExpr) (Value, error) {
	operand, err := interp.evalExpr(e.Operand)
	if err != nil {
		return nil, err
	}
	switch e.Operator {
	case "-":
		switch v := operand.(type) {
		case IntVal:
			return IntVal(-int64(v)), nil
		case FloatVal:
			return FloatVal(-float64(v)), nil
		}
	case "not":
		return BoolVal(!bool(operand.(BoolVal))), nil
	}
	return nil, &RuntimeError{Message: "unsupported unary operation", Line: e.Line}
}

func (interp *Interpreter) evalFuncCall(e *parser.FuncCall) (Value, error) {
	// Built-in: print
	if e.Name == "print" {
		arg, err := interp.evalExpr(e.Args[0])
		if err != nil {
			return nil, err
		}
		interp.io.Print(arg.String())
		return VoidVal{}, nil
	}

	// Built-in: input
	if e.Name == "input" {
		arg, err := interp.evalExpr(e.Args[0])
		if err != nil {
			return nil, err
		}
		result := interp.io.Input(arg.String())
		return StringVal(result), nil
	}

	// User-defined function
	fd, ok := interp.funcs[e.Name]
	if !ok {
		return nil, &RuntimeError{Message: fmt.Sprintf("function '%s' not defined", e.Name), Line: e.Line}
	}

	// Evaluate arguments
	args := make([]Value, len(e.Args))
	for i, arg := range e.Args {
		val, err := interp.evalExpr(arg)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	// Create new scope for function — lexical scoping means parent is global, not caller
	prevEnv := interp.env
	interp.env = NewEnvironment(interp.globalEnv)

	for i, param := range fd.Params {
		interp.env.Define(param.Name, args[i])
	}

	// Execute function body
	var returnVal Value = VoidVal{}
	for _, stmt := range fd.Body {
		err := interp.execStmt(stmt)
		if err != nil {
			if rs, ok := err.(returnSignal); ok {
				returnVal = rs.Value
				break
			}
			interp.env = prevEnv
			return nil, err
		}
	}

	interp.env = prevEnv
	return returnVal, nil
}

func (interp *Interpreter) evalArrayLiteral(e *parser.ArrayLiteral) (Value, error) {
	elements := make([]Value, len(e.Elements))
	for i, elem := range e.Elements {
		val, err := interp.evalExpr(elem)
		if err != nil {
			return nil, err
		}
		elements[i] = val
	}

	var elemType checker.LexoType
	if len(elements) > 0 {
		elemType = elements[0].Type()
	}

	return &ArrayVal{Elements: elements, ElemType: elemType}, nil
}

func (interp *Interpreter) evalIndexAccess(e *parser.IndexAccess) (Value, error) {
	obj, err := interp.evalExpr(e.Object)
	if err != nil {
		return nil, err
	}
	arr := obj.(*ArrayVal)

	idx, err := interp.evalExpr(e.Index)
	if err != nil {
		return nil, err
	}
	i := int(idx.(IntVal))

	if i < 0 || i >= len(arr.Elements) {
		return nil, &RuntimeError{
			Message: fmt.Sprintf("index out of bounds: %d (array length %d)", i, len(arr.Elements)),
			Line:    e.Line,
		}
	}

	return arr.Elements[i], nil
}

func (interp *Interpreter) evalMethodCall(e *parser.MethodCall) (Value, error) {
	obj, err := interp.evalExpr(e.Object)
	if err != nil {
		return nil, err
	}

	switch v := obj.(type) {
	case StringVal:
		s := string(v)
		switch e.Method {
		case "length":
			return IntVal(len(s)), nil
		case "upper":
			return StringVal(strings.ToUpper(s)), nil
		case "lower":
			return StringVal(strings.ToLower(s)), nil
		case "substring":
			startVal, err := interp.evalExpr(e.Args[0])
			if err != nil {
				return nil, err
			}
			endVal, err := interp.evalExpr(e.Args[1])
			if err != nil {
				return nil, err
			}
			start := int(startVal.(IntVal))
			end := int(endVal.(IntVal))
			if start < 0 || end > len(s) || start > end {
				return nil, &RuntimeError{
					Message: fmt.Sprintf("substring indices out of range: [%d:%d] for string of length %d", start, end, len(s)),
					Line:    e.Line,
				}
			}
			return StringVal(s[start:end]), nil
		}
	case *ArrayVal:
		switch e.Method {
		case "length":
			return IntVal(len(v.Elements)), nil
		}
	}

	return nil, &RuntimeError{Message: fmt.Sprintf("unknown method '%s'", e.Method), Line: e.Line}
}
