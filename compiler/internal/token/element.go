package token

type ElementType int

const (
	ElToken ElementType = iota
	ElClass
	ElClassVarDec
	ElType
	ElSubroutineDec
	ElParameterList
	ElSubroutineBody
	ElVarDec
	ElClassName
	ElSubroutineName
	ElVarName
	ElStatements
	ElStatement
	ElLetStatement
	ElIfStatement
	ElWhileStatement
	ElDoStatement
	ElReturnStatement
	ElExpression
	ElTerm
	ElSubroutineCall
	ElExpressionList
	ElOp
	ElUnaryOp
	ElKeywordConstant
)

type Element interface {
	ElementType() ElementType
}
