// Generated from /Users/Nonzero-Sum-Solutions/Google Drive/Projects/Programming/gowork/src/github.com/Nonzero-Sum-Solutions/designlanguage/documentation/design/DesignLanguage.g4 by ANTLR 4.13.1
import org.antlr.v4.runtime.tree.ParseTreeListener;

/**
 * This interface defines a complete listener for a parse tree produced by
 * {@link DesignLanguageParser}.
 */
public interface DesignLanguageListener extends ParseTreeListener {
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#design}.
	 * @param ctx the parse tree
	 */
	void enterDesign(DesignLanguageParser.DesignContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#design}.
	 * @param ctx the parse tree
	 */
	void exitDesign(DesignLanguageParser.DesignContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#preamble}.
	 * @param ctx the parse tree
	 */
	void enterPreamble(DesignLanguageParser.PreambleContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#preamble}.
	 * @param ctx the parse tree
	 */
	void exitPreamble(DesignLanguageParser.PreambleContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#author}.
	 * @param ctx the parse tree
	 */
	void enterAuthor(DesignLanguageParser.AuthorContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#author}.
	 * @param ctx the parse tree
	 */
	void exitAuthor(DesignLanguageParser.AuthorContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#component}.
	 * @param ctx the parse tree
	 */
	void enterComponent(DesignLanguageParser.ComponentContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#component}.
	 * @param ctx the parse tree
	 */
	void exitComponent(DesignLanguageParser.ComponentContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#simpleComponent}.
	 * @param ctx the parse tree
	 */
	void enterSimpleComponent(DesignLanguageParser.SimpleComponentContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#simpleComponent}.
	 * @param ctx the parse tree
	 */
	void exitSimpleComponent(DesignLanguageParser.SimpleComponentContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#field}.
	 * @param ctx the parse tree
	 */
	void enterField(DesignLanguageParser.FieldContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#field}.
	 * @param ctx the parse tree
	 */
	void exitField(DesignLanguageParser.FieldContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#attribute}.
	 * @param ctx the parse tree
	 */
	void enterAttribute(DesignLanguageParser.AttributeContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#attribute}.
	 * @param ctx the parse tree
	 */
	void exitAttribute(DesignLanguageParser.AttributeContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#method}.
	 * @param ctx the parse tree
	 */
	void enterMethod(DesignLanguageParser.MethodContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#method}.
	 * @param ctx the parse tree
	 */
	void exitMethod(DesignLanguageParser.MethodContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#param}.
	 * @param ctx the parse tree
	 */
	void enterParam(DesignLanguageParser.ParamContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#param}.
	 * @param ctx the parse tree
	 */
	void exitParam(DesignLanguageParser.ParamContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#params}.
	 * @param ctx the parse tree
	 */
	void enterParams(DesignLanguageParser.ParamsContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#params}.
	 * @param ctx the parse tree
	 */
	void exitParams(DesignLanguageParser.ParamsContext ctx);
	/**
	 * Enter a parse tree produced by {@link DesignLanguageParser#type}.
	 * @param ctx the parse tree
	 */
	void enterType(DesignLanguageParser.TypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link DesignLanguageParser#type}.
	 * @param ctx the parse tree
	 */
	void exitType(DesignLanguageParser.TypeContext ctx);
}