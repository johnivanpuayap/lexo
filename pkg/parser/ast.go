package parser

// TypeAnnotation represents a type in source code: "int", "string", "int[]", etc.
type TypeAnnotation string

// --- Interfaces ---

type Node interface {
	nodeType() string
	GetLine() int
}

type Statement interface {
	Node
	stmtNode()
}

type Expression interface {
	Node
	exprNode()
}

// --- Program ---

type Program struct {
	Body []Statement
}

// --- Statements ---

type VarDecl struct {
	Name     string
	TypeName TypeAnnotation
	Value    Expression
	Line     int
}

type Assignment struct {
	Target Expression // Identifier or IndexAccess
	Value  Expression
	Line   int
}

type FuncDecl struct {
	Name       string
	Params     []Param
	ReturnType TypeAnnotation
	Body       []Statement
	Line       int
}

type Param struct {
	Name     string
	TypeName TypeAnnotation
}

type ReturnStmt struct {
	Value Expression // nil for bare return in void functions
	Line  int
}

type IfStmt struct {
	Condition Expression
	Body      []Statement
	Elifs     []ElifClause
	ElseBody  []Statement // nil if no else
	Line      int
}

type ElifClause struct {
	Condition Expression
	Body      []Statement
}

type WhileStmt struct {
	Condition Expression
	Body      []Statement
	Line      int
}

type ForStmt struct {
	VarName  string
	VarType  TypeAnnotation
	Iterable Expression
	Body     []Statement
	Line     int
}

type BreakStmt struct {
	Line int
}

type ContinueStmt struct {
	Line int
}

type ExpressionStmt struct {
	Expr Expression
	Line int
}

// --- Expressions ---

type BinaryExpr struct {
	Left     Expression
	Operator string
	Right    Expression
	Line     int
}

type UnaryExpr struct {
	Operator string
	Operand  Expression
	Line     int
}

type Literal struct {
	Value   interface{} // int64, float64, string, or bool
	LitType string      // "int", "float", "string", "bool"
	Line    int
}

type Identifier struct {
	Name string
	Line int
}

type FuncCall struct {
	Name string
	Args []Expression
	Line int
}

type ArrayLiteral struct {
	Elements []Expression
	Line     int
}

type IndexAccess struct {
	Object Expression
	Index  Expression
	Line   int
}

type MethodCall struct {
	Object Expression
	Method string
	Args   []Expression
	Line   int
}

// --- Interface compliance ---

func (s *VarDecl) stmtNode()       {}
func (s *VarDecl) nodeType() string { return "VarDecl" }
func (s *VarDecl) GetLine() int     { return s.Line }

func (s *Assignment) stmtNode()       {}
func (s *Assignment) nodeType() string { return "Assignment" }
func (s *Assignment) GetLine() int     { return s.Line }

func (s *FuncDecl) stmtNode()       {}
func (s *FuncDecl) nodeType() string { return "FuncDecl" }
func (s *FuncDecl) GetLine() int     { return s.Line }

func (s *ReturnStmt) stmtNode()       {}
func (s *ReturnStmt) nodeType() string { return "ReturnStmt" }
func (s *ReturnStmt) GetLine() int     { return s.Line }

func (s *IfStmt) stmtNode()       {}
func (s *IfStmt) nodeType() string { return "IfStmt" }
func (s *IfStmt) GetLine() int     { return s.Line }

func (s *WhileStmt) stmtNode()       {}
func (s *WhileStmt) nodeType() string { return "WhileStmt" }
func (s *WhileStmt) GetLine() int     { return s.Line }

func (s *ForStmt) stmtNode()       {}
func (s *ForStmt) nodeType() string { return "ForStmt" }
func (s *ForStmt) GetLine() int     { return s.Line }

func (s *BreakStmt) stmtNode()       {}
func (s *BreakStmt) nodeType() string { return "BreakStmt" }
func (s *BreakStmt) GetLine() int     { return s.Line }

func (s *ContinueStmt) stmtNode()       {}
func (s *ContinueStmt) nodeType() string { return "ContinueStmt" }
func (s *ContinueStmt) GetLine() int     { return s.Line }

func (s *ExpressionStmt) stmtNode()       {}
func (s *ExpressionStmt) nodeType() string { return "ExpressionStmt" }
func (s *ExpressionStmt) GetLine() int     { return s.Line }

func (e *BinaryExpr) exprNode()       {}
func (e *BinaryExpr) nodeType() string { return "BinaryExpr" }
func (e *BinaryExpr) GetLine() int     { return e.Line }

func (e *UnaryExpr) exprNode()       {}
func (e *UnaryExpr) nodeType() string { return "UnaryExpr" }
func (e *UnaryExpr) GetLine() int     { return e.Line }

func (e *Literal) exprNode()       {}
func (e *Literal) nodeType() string { return "Literal" }
func (e *Literal) GetLine() int     { return e.Line }

func (e *Identifier) exprNode()       {}
func (e *Identifier) nodeType() string { return "Identifier" }
func (e *Identifier) GetLine() int     { return e.Line }

func (e *FuncCall) exprNode()       {}
func (e *FuncCall) nodeType() string { return "FuncCall" }
func (e *FuncCall) GetLine() int     { return e.Line }

func (e *ArrayLiteral) exprNode()       {}
func (e *ArrayLiteral) nodeType() string { return "ArrayLiteral" }
func (e *ArrayLiteral) GetLine() int     { return e.Line }

func (e *IndexAccess) exprNode()       {}
func (e *IndexAccess) nodeType() string { return "IndexAccess" }
func (e *IndexAccess) GetLine() int     { return e.Line }

func (e *MethodCall) exprNode()       {}
func (e *MethodCall) nodeType() string { return "MethodCall" }
func (e *MethodCall) GetLine() int     { return e.Line }
