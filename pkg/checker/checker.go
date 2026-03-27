package checker

import (
	"fmt"

	"github.com/johnivanpuayap/lexo/pkg/parser"
)

type TypeError struct {
	Message string
	Line    int
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("Line %d: %s", e.Line, e.Message)
}

type FuncSig struct {
	Params     []LexoType
	ReturnType LexoType
}

type scope struct {
	vars   map[string]LexoType
	parent *scope
}

func newScope(parent *scope) *scope {
	return &scope{vars: make(map[string]LexoType), parent: parent}
}

func (s *scope) get(name string) (LexoType, bool) {
	if t, ok := s.vars[name]; ok {
		return t, true
	}
	if s.parent != nil {
		return s.parent.get(name)
	}
	return 0, false
}

func (s *scope) set(name string, t LexoType) {
	s.vars[name] = t
}

func (s *scope) has(name string) bool {
	_, ok := s.vars[name]
	return ok
}

type checker struct {
	scope      *scope
	funcs      map[string]FuncSig
	returnType *LexoType // nil when not inside a function
	loopDepth  int
}

func Check(prog *parser.Program) error {
	c := &checker{
		scope: newScope(nil),
		funcs: map[string]FuncSig{
			"print": {Params: nil, ReturnType: TypeVoid},   // special: accepts any single arg
			"input": {Params: []LexoType{TypeString}, ReturnType: TypeString},
		},
	}

	// First pass: register all function signatures
	for _, stmt := range prog.Body {
		if fd, ok := stmt.(*parser.FuncDecl); ok {
			sig, err := c.buildFuncSig(fd)
			if err != nil {
				return err
			}
			c.funcs[fd.Name] = sig
		}
	}

	// Second pass: check all statements
	for _, stmt := range prog.Body {
		if err := c.checkStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) buildFuncSig(fd *parser.FuncDecl) (FuncSig, error) {
	var params []LexoType
	for _, p := range fd.Params {
		t, ok := TypeFromAnnotation(string(p.TypeName))
		if !ok {
			return FuncSig{}, &TypeError{Message: fmt.Sprintf("unknown type '%s'", p.TypeName), Line: fd.Line}
		}
		params = append(params, t)
	}
	rt, ok := TypeFromAnnotation(string(fd.ReturnType))
	if !ok {
		return FuncSig{}, &TypeError{Message: fmt.Sprintf("unknown return type '%s'", fd.ReturnType), Line: fd.Line}
	}
	return FuncSig{Params: params, ReturnType: rt}, nil
}

func (c *checker) checkStmt(stmt parser.Statement) error {
	switch s := stmt.(type) {
	case *parser.VarDecl:
		return c.checkVarDecl(s)
	case *parser.Assignment:
		return c.checkAssignment(s)
	case *parser.FuncDecl:
		return c.checkFuncDecl(s)
	case *parser.ReturnStmt:
		return c.checkReturnStmt(s)
	case *parser.IfStmt:
		return c.checkIfStmt(s)
	case *parser.WhileStmt:
		return c.checkWhileStmt(s)
	case *parser.ForStmt:
		return c.checkForStmt(s)
	case *parser.BreakStmt:
		if c.loopDepth == 0 {
			return &TypeError{Message: "'break' outside of loop", Line: s.Line}
		}
		return nil
	case *parser.ContinueStmt:
		if c.loopDepth == 0 {
			return &TypeError{Message: "'continue' outside of loop", Line: s.Line}
		}
		return nil
	case *parser.ExpressionStmt:
		_, err := c.checkExpr(s.Expr)
		return err
	default:
		return &TypeError{Message: fmt.Sprintf("unknown statement type %T", stmt), Line: stmt.GetLine()}
	}
}

func (c *checker) checkVarDecl(s *parser.VarDecl) error {
	if c.scope.has(s.Name) {
		return &TypeError{
			Message: fmt.Sprintf("variable '%s' is already declared in this scope", s.Name),
			Line:    s.Line,
		}
	}

	declType, ok := TypeFromAnnotation(string(s.TypeName))
	if !ok {
		return &TypeError{Message: fmt.Sprintf("unknown type '%s'", s.TypeName), Line: s.Line}
	}

	valType, err := c.checkExpr(s.Value)
	if err != nil {
		return err
	}

	if valType != declType {
		return &TypeError{
			Message: fmt.Sprintf("type mismatch: cannot assign %s value to %s variable '%s'", valType, declType, s.Name),
			Line:    s.Line,
		}
	}

	c.scope.set(s.Name, declType)
	return nil
}

func (c *checker) checkAssignment(s *parser.Assignment) error {
	valType, err := c.checkExpr(s.Value)
	if err != nil {
		return err
	}

	switch target := s.Target.(type) {
	case *parser.Identifier:
		varType, ok := c.scope.get(target.Name)
		if !ok {
			return &TypeError{
				Message: fmt.Sprintf("variable '%s' is not declared", target.Name),
				Line:    s.Line,
			}
		}
		if valType != varType {
			return &TypeError{
				Message: fmt.Sprintf("type mismatch: cannot assign %s value to %s variable '%s'", valType, varType, target.Name),
				Line:    s.Line,
			}
		}
	case *parser.IndexAccess:
		arrType, err := c.checkExpr(target.Object)
		if err != nil {
			return err
		}
		elemType, ok := ElementType(arrType)
		if !ok {
			return &TypeError{Message: "cannot index into non-array type", Line: s.Line}
		}
		idxType, err := c.checkExpr(target.Index)
		if err != nil {
			return err
		}
		if idxType != TypeInt {
			return &TypeError{Message: "array index must be an int", Line: s.Line}
		}
		if valType != elemType {
			return &TypeError{
				Message: fmt.Sprintf("type mismatch: cannot assign %s to %s array element", valType, elemType),
				Line:    s.Line,
			}
		}
	default:
		return &TypeError{Message: "invalid assignment target", Line: s.Line}
	}
	return nil
}

func (c *checker) checkFuncDecl(fd *parser.FuncDecl) error {
	sig := c.funcs[fd.Name]

	prevScope := c.scope
	c.scope = newScope(c.scope)

	for i, p := range fd.Params {
		c.scope.set(p.Name, sig.Params[i])
	}

	prevReturn := c.returnType
	c.returnType = &sig.ReturnType

	for _, stmt := range fd.Body {
		if err := c.checkStmt(stmt); err != nil {
			return err
		}
	}

	c.returnType = prevReturn
	c.scope = prevScope
	return nil
}

func (c *checker) checkReturnStmt(s *parser.ReturnStmt) error {
	if c.returnType == nil {
		return &TypeError{Message: "return statement outside of function", Line: s.Line}
	}

	if s.Value == nil {
		if *c.returnType != TypeVoid {
			return &TypeError{
				Message: fmt.Sprintf("return type mismatch: function expects %s but got no return value", *c.returnType),
				Line:    s.Line,
			}
		}
		return nil
	}

	valType, err := c.checkExpr(s.Value)
	if err != nil {
		return err
	}

	if valType != *c.returnType {
		return &TypeError{
			Message: fmt.Sprintf("return type mismatch: function expects %s but got %s", *c.returnType, valType),
			Line:    s.Line,
		}
	}
	return nil
}

func (c *checker) checkIfStmt(s *parser.IfStmt) error {
	condType, err := c.checkExpr(s.Condition)
	if err != nil {
		return err
	}
	if condType != TypeBool {
		return &TypeError{
			Message: fmt.Sprintf("condition must be bool, got %s", condType),
			Line:    s.Line,
		}
	}

	for _, stmt := range s.Body {
		if err := c.checkStmt(stmt); err != nil {
			return err
		}
	}

	for _, elif := range s.Elifs {
		ct, err := c.checkExpr(elif.Condition)
		if err != nil {
			return err
		}
		if ct != TypeBool {
			return &TypeError{Message: fmt.Sprintf("elif condition must be bool, got %s", ct), Line: s.Line}
		}
		for _, stmt := range elif.Body {
			if err := c.checkStmt(stmt); err != nil {
				return err
			}
		}
	}

	for _, stmt := range s.ElseBody {
		if err := c.checkStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) checkWhileStmt(s *parser.WhileStmt) error {
	condType, err := c.checkExpr(s.Condition)
	if err != nil {
		return err
	}
	if condType != TypeBool {
		return &TypeError{Message: fmt.Sprintf("while condition must be bool, got %s", condType), Line: s.Line}
	}
	c.loopDepth++
	for _, stmt := range s.Body {
		if err := c.checkStmt(stmt); err != nil {
			c.loopDepth--
			return err
		}
	}
	c.loopDepth--
	return nil
}

func (c *checker) checkForStmt(s *parser.ForStmt) error {
	iterType, err := c.checkExpr(s.Iterable)
	if err != nil {
		return err
	}
	elemType, ok := ElementType(iterType)
	if !ok {
		return &TypeError{Message: fmt.Sprintf("cannot iterate over %s, expected an array", iterType), Line: s.Line}
	}

	varType, ok := TypeFromAnnotation(string(s.VarType))
	if !ok {
		return &TypeError{Message: fmt.Sprintf("unknown type '%s'", s.VarType), Line: s.Line}
	}

	if varType != elemType {
		return &TypeError{
			Message: fmt.Sprintf("for loop variable type %s does not match array element type %s", varType, elemType),
			Line:    s.Line,
		}
	}

	prevScope := c.scope
	c.scope = newScope(c.scope)
	c.scope.set(s.VarName, varType)

	c.loopDepth++
	for _, stmt := range s.Body {
		if err := c.checkStmt(stmt); err != nil {
			c.loopDepth--
			c.scope = prevScope
			return err
		}
	}
	c.loopDepth--

	c.scope = prevScope
	return nil
}

func (c *checker) checkExpr(expr parser.Expression) (LexoType, error) {
	switch e := expr.(type) {
	case *parser.Literal:
		return c.checkLiteral(e)
	case *parser.Identifier:
		t, ok := c.scope.get(e.Name)
		if !ok {
			return 0, &TypeError{Message: fmt.Sprintf("variable '%s' is not declared", e.Name), Line: e.Line}
		}
		return t, nil
	case *parser.BinaryExpr:
		return c.checkBinaryExpr(e)
	case *parser.UnaryExpr:
		return c.checkUnaryExpr(e)
	case *parser.FuncCall:
		return c.checkFuncCall(e)
	case *parser.ArrayLiteral:
		return c.checkArrayLiteral(e)
	case *parser.IndexAccess:
		return c.checkIndexAccess(e)
	case *parser.MethodCall:
		return c.checkMethodCall(e)
	default:
		return 0, &TypeError{Message: fmt.Sprintf("unknown expression type %T", expr), Line: expr.GetLine()}
	}
}

func (c *checker) checkLiteral(e *parser.Literal) (LexoType, error) {
	switch e.LitType {
	case "int":
		return TypeInt, nil
	case "float":
		return TypeFloat, nil
	case "string":
		return TypeString, nil
	case "bool":
		return TypeBool, nil
	default:
		return 0, &TypeError{Message: fmt.Sprintf("unknown literal type '%s'", e.LitType), Line: e.Line}
	}
}

func (c *checker) checkBinaryExpr(e *parser.BinaryExpr) (LexoType, error) {
	leftType, err := c.checkExpr(e.Left)
	if err != nil {
		return 0, err
	}
	rightType, err := c.checkExpr(e.Right)
	if err != nil {
		return 0, err
	}

	switch e.Operator {
	case "+":
		if leftType == TypeString && rightType == TypeString {
			return TypeString, nil
		}
		if leftType != rightType {
			return 0, &TypeError{
				Message: fmt.Sprintf("type mismatch in '+': cannot add %s and %s", leftType, rightType),
				Line:    e.Line,
			}
		}
		if leftType == TypeInt || leftType == TypeFloat {
			return leftType, nil
		}
		return 0, &TypeError{Message: fmt.Sprintf("operator '+' not supported for %s", leftType), Line: e.Line}

	case "-", "*", "/", "%":
		if leftType != rightType {
			return 0, &TypeError{
				Message: fmt.Sprintf("type mismatch in '%s': cannot operate on %s and %s", e.Operator, leftType, rightType),
				Line:    e.Line,
			}
		}
		if leftType == TypeInt || leftType == TypeFloat {
			return leftType, nil
		}
		return 0, &TypeError{
			Message: fmt.Sprintf("operator '%s' not supported for %s", e.Operator, leftType),
			Line:    e.Line,
		}

	case "==", "!=":
		if leftType != rightType {
			return 0, &TypeError{
				Message: fmt.Sprintf("type mismatch in '%s': cannot compare %s and %s", e.Operator, leftType, rightType),
				Line:    e.Line,
			}
		}
		return TypeBool, nil

	case "<", ">", "<=", ">=":
		if leftType != rightType {
			return 0, &TypeError{
				Message: fmt.Sprintf("type mismatch in '%s': cannot compare %s and %s", e.Operator, leftType, rightType),
				Line:    e.Line,
			}
		}
		if leftType == TypeInt || leftType == TypeFloat {
			return TypeBool, nil
		}
		return 0, &TypeError{
			Message: fmt.Sprintf("operator '%s' not supported for %s", e.Operator, leftType),
			Line:    e.Line,
		}

	case "and", "or":
		if leftType != TypeBool || rightType != TypeBool {
			return 0, &TypeError{
				Message: fmt.Sprintf("operator '%s' requires bool operands, got %s and %s", e.Operator, leftType, rightType),
				Line:    e.Line,
			}
		}
		return TypeBool, nil

	default:
		return 0, &TypeError{Message: fmt.Sprintf("unknown operator '%s'", e.Operator), Line: e.Line}
	}
}

func (c *checker) checkUnaryExpr(e *parser.UnaryExpr) (LexoType, error) {
	operandType, err := c.checkExpr(e.Operand)
	if err != nil {
		return 0, err
	}

	switch e.Operator {
	case "-":
		if operandType == TypeInt || operandType == TypeFloat {
			return operandType, nil
		}
		return 0, &TypeError{Message: fmt.Sprintf("unary '-' not supported for %s", operandType), Line: e.Line}
	case "not":
		if operandType == TypeBool {
			return TypeBool, nil
		}
		return 0, &TypeError{Message: fmt.Sprintf("'not' requires bool, got %s", operandType), Line: e.Line}
	default:
		return 0, &TypeError{Message: fmt.Sprintf("unknown unary operator '%s'", e.Operator), Line: e.Line}
	}
}

func (c *checker) checkFuncCall(e *parser.FuncCall) (LexoType, error) {
	// Special handling for print (accepts any single argument)
	if e.Name == "print" {
		if len(e.Args) != 1 {
			return 0, &TypeError{
				Message: fmt.Sprintf("print() takes 1 argument, got %d", len(e.Args)),
				Line:    e.Line,
			}
		}
		_, err := c.checkExpr(e.Args[0])
		if err != nil {
			return 0, err
		}
		return TypeVoid, nil
	}

	sig, ok := c.funcs[e.Name]
	if !ok {
		return 0, &TypeError{Message: fmt.Sprintf("function '%s' is not defined", e.Name), Line: e.Line}
	}

	if len(e.Args) != len(sig.Params) {
		return 0, &TypeError{
			Message: fmt.Sprintf("function '%s' takes %d argument(s), got %d", e.Name, len(sig.Params), len(e.Args)),
			Line:    e.Line,
		}
	}

	for i, arg := range e.Args {
		argType, err := c.checkExpr(arg)
		if err != nil {
			return 0, err
		}
		if argType != sig.Params[i] {
			return 0, &TypeError{
				Message: fmt.Sprintf("argument %d of '%s': expected %s, got %s", i+1, e.Name, sig.Params[i], argType),
				Line:    e.Line,
			}
		}
	}

	return sig.ReturnType, nil
}

func (c *checker) checkArrayLiteral(e *parser.ArrayLiteral) (LexoType, error) {
	if len(e.Elements) == 0 {
		return 0, &TypeError{Message: "cannot infer type of empty array literal", Line: e.Line}
	}

	firstType, err := c.checkExpr(e.Elements[0])
	if err != nil {
		return 0, err
	}

	for i := 1; i < len(e.Elements); i++ {
		elemType, err := c.checkExpr(e.Elements[i])
		if err != nil {
			return 0, err
		}
		if elemType != firstType {
			return 0, &TypeError{
				Message: fmt.Sprintf("array elements must be the same type: expected %s, got %s at position %d", firstType, elemType, i),
				Line:    e.Line,
			}
		}
	}

	arrType, ok := ArrayTypeOf(firstType)
	if !ok {
		return 0, &TypeError{Message: fmt.Sprintf("cannot create array of %s", firstType), Line: e.Line}
	}
	return arrType, nil
}

func (c *checker) checkIndexAccess(e *parser.IndexAccess) (LexoType, error) {
	objType, err := c.checkExpr(e.Object)
	if err != nil {
		return 0, err
	}

	elemType, ok := ElementType(objType)
	if !ok {
		return 0, &TypeError{Message: fmt.Sprintf("cannot index into %s — expected an array", objType), Line: e.Line}
	}

	idxType, err := c.checkExpr(e.Index)
	if err != nil {
		return 0, err
	}
	if idxType != TypeInt {
		return 0, &TypeError{Message: fmt.Sprintf("array index must be int, got %s", idxType), Line: e.Line}
	}

	return elemType, nil
}

func (c *checker) checkMethodCall(e *parser.MethodCall) (LexoType, error) {
	objType, err := c.checkExpr(e.Object)
	if err != nil {
		return 0, err
	}

	switch {
	case objType == TypeString:
		switch e.Method {
		case "length":
			if len(e.Args) != 0 {
				return 0, &TypeError{Message: "length() takes no arguments", Line: e.Line}
			}
			return TypeInt, nil
		case "upper", "lower":
			if len(e.Args) != 0 {
				return 0, &TypeError{Message: fmt.Sprintf("%s() takes no arguments", e.Method), Line: e.Line}
			}
			return TypeString, nil
		case "substring":
			if len(e.Args) != 2 {
				return 0, &TypeError{Message: "substring() takes 2 arguments (start, end)", Line: e.Line}
			}
			for i, arg := range e.Args {
				at, err := c.checkExpr(arg)
				if err != nil {
					return 0, err
				}
				if at != TypeInt {
					return 0, &TypeError{
						Message: fmt.Sprintf("substring() argument %d must be int, got %s", i+1, at),
						Line:    e.Line,
					}
				}
			}
			return TypeString, nil
		default:
			return 0, &TypeError{Message: fmt.Sprintf("string has no method '%s'", e.Method), Line: e.Line}
		}

	case IsArrayType(objType):
		switch e.Method {
		case "length":
			if len(e.Args) != 0 {
				return 0, &TypeError{Message: "length() takes no arguments", Line: e.Line}
			}
			return TypeInt, nil
		default:
			return 0, &TypeError{Message: fmt.Sprintf("%s has no method '%s'", objType, e.Method), Line: e.Line}
		}

	default:
		return 0, &TypeError{Message: fmt.Sprintf("type %s has no methods", objType), Line: e.Line}
	}
}
