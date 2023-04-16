package tokenizer

//go:generate go run golang.org/x/tools/cmd/stringer -type=ElementType -linecomment
type ElementType int

const (
	ElToken           ElementType = iota //token
	ElClass                              //class
	ElClassVarDec                        //classVarDec
	ElType                               //type
	ElSubroutineDec                      //subroutineDec
	ElParameterList                      //parameterList
	ElSubroutineBody                     //subroutineBody
	ElVarDec                             //varDec
	ElClassName                          //className
	ElSubroutineName                     //subroutineName
	ElVarName                            //varName
	ElStatements                         //statements
	ElStatement                          //statement
	ElLetStatement                       //letStatement
	ElIfStatement                        //ifStatement
	ElWhileStatement                     //whileStatement
	ElDoStatement                        //doStatement
	ElReturnStatement                    //returnStatement
	ElExpression                         //expression
	ElTerm                               //term
	ElSubroutineCall                     //subroutineCall
	ElExpressionList                     //expressionList
	ElOp                                 //op
	ElUnaryOp                            //unaryOp
	ElKeywordConstant                    //keywordConstant
)

type Element interface {
	ElementType() ElementType
}
