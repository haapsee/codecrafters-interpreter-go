package interfaces

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

type Visitor interface {
	VisitBinaryExpr(b Expr) (interface{}, error)
	VisitGroupingExpr(g Expr) (interface{}, error)
	VisitLiteralExpr(l Expr) (interface{}, error)
	VisitUnaryExpr(u Expr) (interface{}, error)
}
