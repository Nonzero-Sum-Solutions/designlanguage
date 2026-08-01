// Generated from /Users/Nonzero-Sum-Solutions/Google Drive/Projects/Programming/gowork/src/github.com/Nonzero-Sum-Solutions/designlanguage/documentation/design/DesignLanguage.g4 by ANTLR 4.9.2
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast"})
public class DesignLanguageParser extends Parser {
	static { RuntimeMetaData.checkVersion("4.9.2", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		T__0=1, T__1=2, T__2=3, T__3=4, ARRAY=5, ARROW=6, ASTERISK=7, FIELD_START=8, 
		AUTHOR_NAME=9, AUTHOR_START=10, COMMENT=11, COMMENT_START=12, COMMENT_TEXT=13, 
		NAME=14, NEWLINE=15, SPECIAL_CHAR=16;
	public static final int
		RULE_design = 0, RULE_preamble = 1, RULE_author = 2, RULE_component = 3, 
		RULE_simpleComponent = 4, RULE_field = 5, RULE_attribute = 6, RULE_method = 7, 
		RULE_param = 8, RULE_params = 9, RULE_type = 10;
	private static String[] makeRuleNames() {
		return new String[] {
			"design", "preamble", "author", "component", "simpleComponent", "field", 
			"attribute", "method", "param", "params", "type"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "' '", "'()'", "'('", "')'", "'[]'", "'->'", "'*'", null, null, 
			null, null, "'-- '"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, null, null, null, null, "ARRAY", "ARROW", "ASTERISK", "FIELD_START", 
			"AUTHOR_NAME", "AUTHOR_START", "COMMENT", "COMMENT_START", "COMMENT_TEXT", 
			"NAME", "NEWLINE", "SPECIAL_CHAR"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}

	@Override
	public String getGrammarFileName() { return "DesignLanguage.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public DesignLanguageParser(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	public static class DesignContext extends ParserRuleContext {
		public TerminalNode EOF() { return getToken(DesignLanguageParser.EOF, 0); }
		public PreambleContext preamble() {
			return getRuleContext(PreambleContext.class,0);
		}
		public List<ComponentContext> component() {
			return getRuleContexts(ComponentContext.class);
		}
		public ComponentContext component(int i) {
			return getRuleContext(ComponentContext.class,i);
		}
		public List<TerminalNode> NEWLINE() { return getTokens(DesignLanguageParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(DesignLanguageParser.NEWLINE, i);
		}
		public DesignContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_design; }
	}

	public final DesignContext design() throws RecognitionException {
		DesignContext _localctx = new DesignContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_design);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(23);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==AUTHOR_START || _la==COMMENT) {
				{
				setState(22);
				preamble();
				}
			}

			setState(30);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==NAME) {
				{
				{
				setState(25);
				component();
				setState(26);
				match(NEWLINE);
				}
				}
				setState(32);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(33);
			match(EOF);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class PreambleContext extends ParserRuleContext {
		public List<TerminalNode> NEWLINE() { return getTokens(DesignLanguageParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(DesignLanguageParser.NEWLINE, i);
		}
		public AuthorContext author() {
			return getRuleContext(AuthorContext.class,0);
		}
		public TerminalNode COMMENT() { return getToken(DesignLanguageParser.COMMENT, 0); }
		public PreambleContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_preamble; }
	}

	public final PreambleContext preamble() throws RecognitionException {
		PreambleContext _localctx = new PreambleContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_preamble);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(43);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case AUTHOR_START:
				{
				setState(35);
				author();
				setState(36);
				match(NEWLINE);
				setState(39);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==COMMENT) {
					{
					setState(37);
					match(COMMENT);
					setState(38);
					match(NEWLINE);
					}
				}

				}
				break;
			case COMMENT:
				{
				{
				setState(41);
				match(COMMENT);
				setState(42);
				match(NEWLINE);
				}
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			setState(45);
			match(NEWLINE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class AuthorContext extends ParserRuleContext {
		public TerminalNode AUTHOR_START() { return getToken(DesignLanguageParser.AUTHOR_START, 0); }
		public TerminalNode AUTHOR_NAME() { return getToken(DesignLanguageParser.AUTHOR_NAME, 0); }
		public AuthorContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_author; }
	}

	public final AuthorContext author() throws RecognitionException {
		AuthorContext _localctx = new AuthorContext(_ctx, getState());
		enterRule(_localctx, 4, RULE_author);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(47);
			match(AUTHOR_START);
			setState(48);
			match(AUTHOR_NAME);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ComponentContext extends ParserRuleContext {
		public SimpleComponentContext simpleComponent() {
			return getRuleContext(SimpleComponentContext.class,0);
		}
		public List<FieldContext> field() {
			return getRuleContexts(FieldContext.class);
		}
		public FieldContext field(int i) {
			return getRuleContext(FieldContext.class,i);
		}
		public List<TerminalNode> NEWLINE() { return getTokens(DesignLanguageParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(DesignLanguageParser.NEWLINE, i);
		}
		public ComponentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_component; }
	}

	public final ComponentContext component() throws RecognitionException {
		ComponentContext _localctx = new ComponentContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_component);
		int _la;
		try {
			setState(59);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,5,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(50);
				simpleComponent();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(51);
				simpleComponent();
				setState(55); 
				_errHandler.sync(this);
				_la = _input.LA(1);
				do {
					{
					{
					setState(52);
					field();
					setState(53);
					match(NEWLINE);
					}
					}
					setState(57); 
					_errHandler.sync(this);
					_la = _input.LA(1);
				} while ( _la==FIELD_START );
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SimpleComponentContext extends ParserRuleContext {
		public TerminalNode NAME() { return getToken(DesignLanguageParser.NAME, 0); }
		public List<TerminalNode> NEWLINE() { return getTokens(DesignLanguageParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(DesignLanguageParser.NEWLINE, i);
		}
		public TerminalNode COMMENT() { return getToken(DesignLanguageParser.COMMENT, 0); }
		public SimpleComponentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_simpleComponent; }
	}

	public final SimpleComponentContext simpleComponent() throws RecognitionException {
		SimpleComponentContext _localctx = new SimpleComponentContext(_ctx, getState());
		enterRule(_localctx, 8, RULE_simpleComponent);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(61);
			match(NAME);
			setState(62);
			match(NEWLINE);
			setState(65);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMENT) {
				{
				setState(63);
				match(COMMENT);
				setState(64);
				match(NEWLINE);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class FieldContext extends ParserRuleContext {
		public TerminalNode FIELD_START() { return getToken(DesignLanguageParser.FIELD_START, 0); }
		public AttributeContext attribute() {
			return getRuleContext(AttributeContext.class,0);
		}
		public MethodContext method() {
			return getRuleContext(MethodContext.class,0);
		}
		public FieldContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_field; }
	}

	public final FieldContext field() throws RecognitionException {
		FieldContext _localctx = new FieldContext(_ctx, getState());
		enterRule(_localctx, 10, RULE_field);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(67);
			match(FIELD_START);
			setState(70);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,7,_ctx) ) {
			case 1:
				{
				setState(68);
				attribute();
				}
				break;
			case 2:
				{
				setState(69);
				method();
				}
				break;
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class AttributeContext extends ParserRuleContext {
		public ParamContext param() {
			return getRuleContext(ParamContext.class,0);
		}
		public TerminalNode COMMENT() { return getToken(DesignLanguageParser.COMMENT, 0); }
		public AttributeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_attribute; }
	}

	public final AttributeContext attribute() throws RecognitionException {
		AttributeContext _localctx = new AttributeContext(_ctx, getState());
		enterRule(_localctx, 12, RULE_attribute);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(72);
			param();
			setState(75);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==T__0) {
				{
				setState(73);
				match(T__0);
				setState(74);
				match(COMMENT);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class MethodContext extends ParserRuleContext {
		public TerminalNode NAME() { return getToken(DesignLanguageParser.NAME, 0); }
		public List<ParamsContext> params() {
			return getRuleContexts(ParamsContext.class);
		}
		public ParamsContext params(int i) {
			return getRuleContext(ParamsContext.class,i);
		}
		public TerminalNode ARROW() { return getToken(DesignLanguageParser.ARROW, 0); }
		public TerminalNode COMMENT() { return getToken(DesignLanguageParser.COMMENT, 0); }
		public MethodContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_method; }
	}

	public final MethodContext method() throws RecognitionException {
		MethodContext _localctx = new MethodContext(_ctx, getState());
		enterRule(_localctx, 14, RULE_method);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(77);
			match(NAME);
			setState(78);
			match(T__0);
			setState(79);
			params();
			setState(80);
			match(T__0);
			setState(84);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ARROW) {
				{
				setState(81);
				match(ARROW);
				setState(82);
				match(T__0);
				setState(83);
				params();
				}
			}

			setState(88);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==T__0) {
				{
				setState(86);
				match(T__0);
				setState(87);
				match(COMMENT);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ParamContext extends ParserRuleContext {
		public TerminalNode NAME() { return getToken(DesignLanguageParser.NAME, 0); }
		public TypeContext type() {
			return getRuleContext(TypeContext.class,0);
		}
		public ParamContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_param; }
	}

	public final ParamContext param() throws RecognitionException {
		ParamContext _localctx = new ParamContext(_ctx, getState());
		enterRule(_localctx, 16, RULE_param);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(90);
			match(NAME);
			setState(91);
			match(T__0);
			setState(92);
			type();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ParamsContext extends ParserRuleContext {
		public List<ParamContext> param() {
			return getRuleContexts(ParamContext.class);
		}
		public ParamContext param(int i) {
			return getRuleContext(ParamContext.class,i);
		}
		public ParamsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_params; }
	}

	public final ParamsContext params() throws RecognitionException {
		ParamsContext _localctx = new ParamsContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_params);
		int _la;
		try {
			setState(107);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case T__1:
				enterOuterAlt(_localctx, 1);
				{
				setState(94);
				match(T__1);
				}
				break;
			case T__2:
				enterOuterAlt(_localctx, 2);
				{
				setState(95);
				match(T__2);
				setState(97);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==NAME) {
					{
					setState(96);
					param();
					}
				}

				setState(103);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while (_la==T__0) {
					{
					{
					setState(99);
					match(T__0);
					setState(100);
					param();
					}
					}
					setState(105);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				setState(106);
				match(T__3);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class TypeContext extends ParserRuleContext {
		public TerminalNode NAME() { return getToken(DesignLanguageParser.NAME, 0); }
		public TerminalNode ARRAY() { return getToken(DesignLanguageParser.ARRAY, 0); }
		public TypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_type; }
	}

	public final TypeContext type() throws RecognitionException {
		TypeContext _localctx = new TypeContext(_ctx, getState());
		enterRule(_localctx, 20, RULE_type);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(110);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ARRAY) {
				{
				setState(109);
				match(ARRAY);
				}
			}

			setState(112);
			match(NAME);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static final String _serializedATN =
		"\3\u608b\ua72a\u8133\ub9ed\u417c\u3be7\u7786\u5964\3\22u\4\2\t\2\4\3\t"+
		"\3\4\4\t\4\4\5\t\5\4\6\t\6\4\7\t\7\4\b\t\b\4\t\t\t\4\n\t\n\4\13\t\13\4"+
		"\f\t\f\3\2\5\2\32\n\2\3\2\3\2\3\2\7\2\37\n\2\f\2\16\2\"\13\2\3\2\3\2\3"+
		"\3\3\3\3\3\3\3\5\3*\n\3\3\3\3\3\5\3.\n\3\3\3\3\3\3\4\3\4\3\4\3\5\3\5\3"+
		"\5\3\5\3\5\6\5:\n\5\r\5\16\5;\5\5>\n\5\3\6\3\6\3\6\3\6\5\6D\n\6\3\7\3"+
		"\7\3\7\5\7I\n\7\3\b\3\b\3\b\5\bN\n\b\3\t\3\t\3\t\3\t\3\t\3\t\3\t\5\tW"+
		"\n\t\3\t\3\t\5\t[\n\t\3\n\3\n\3\n\3\n\3\13\3\13\3\13\5\13d\n\13\3\13\3"+
		"\13\7\13h\n\13\f\13\16\13k\13\13\3\13\5\13n\n\13\3\f\5\fq\n\f\3\f\3\f"+
		"\3\f\2\2\r\2\4\6\b\n\f\16\20\22\24\26\2\2\2x\2\31\3\2\2\2\4-\3\2\2\2\6"+
		"\61\3\2\2\2\b=\3\2\2\2\n?\3\2\2\2\fE\3\2\2\2\16J\3\2\2\2\20O\3\2\2\2\22"+
		"\\\3\2\2\2\24m\3\2\2\2\26p\3\2\2\2\30\32\5\4\3\2\31\30\3\2\2\2\31\32\3"+
		"\2\2\2\32 \3\2\2\2\33\34\5\b\5\2\34\35\7\21\2\2\35\37\3\2\2\2\36\33\3"+
		"\2\2\2\37\"\3\2\2\2 \36\3\2\2\2 !\3\2\2\2!#\3\2\2\2\" \3\2\2\2#$\7\2\2"+
		"\3$\3\3\2\2\2%&\5\6\4\2&)\7\21\2\2\'(\7\r\2\2(*\7\21\2\2)\'\3\2\2\2)*"+
		"\3\2\2\2*.\3\2\2\2+,\7\r\2\2,.\7\21\2\2-%\3\2\2\2-+\3\2\2\2./\3\2\2\2"+
		"/\60\7\21\2\2\60\5\3\2\2\2\61\62\7\f\2\2\62\63\7\13\2\2\63\7\3\2\2\2\64"+
		">\5\n\6\2\659\5\n\6\2\66\67\5\f\7\2\678\7\21\2\28:\3\2\2\29\66\3\2\2\2"+
		":;\3\2\2\2;9\3\2\2\2;<\3\2\2\2<>\3\2\2\2=\64\3\2\2\2=\65\3\2\2\2>\t\3"+
		"\2\2\2?@\7\20\2\2@C\7\21\2\2AB\7\r\2\2BD\7\21\2\2CA\3\2\2\2CD\3\2\2\2"+
		"D\13\3\2\2\2EH\7\n\2\2FI\5\16\b\2GI\5\20\t\2HF\3\2\2\2HG\3\2\2\2I\r\3"+
		"\2\2\2JM\5\22\n\2KL\7\3\2\2LN\7\r\2\2MK\3\2\2\2MN\3\2\2\2N\17\3\2\2\2"+
		"OP\7\20\2\2PQ\7\3\2\2QR\5\24\13\2RV\7\3\2\2ST\7\b\2\2TU\7\3\2\2UW\5\24"+
		"\13\2VS\3\2\2\2VW\3\2\2\2WZ\3\2\2\2XY\7\3\2\2Y[\7\r\2\2ZX\3\2\2\2Z[\3"+
		"\2\2\2[\21\3\2\2\2\\]\7\20\2\2]^\7\3\2\2^_\5\26\f\2_\23\3\2\2\2`n\7\4"+
		"\2\2ac\7\5\2\2bd\5\22\n\2cb\3\2\2\2cd\3\2\2\2di\3\2\2\2ef\7\3\2\2fh\5"+
		"\22\n\2ge\3\2\2\2hk\3\2\2\2ig\3\2\2\2ij\3\2\2\2jl\3\2\2\2ki\3\2\2\2ln"+
		"\7\6\2\2m`\3\2\2\2ma\3\2\2\2n\25\3\2\2\2oq\7\7\2\2po\3\2\2\2pq\3\2\2\2"+
		"qr\3\2\2\2rs\7\20\2\2s\27\3\2\2\2\21\31 )-;=CHMVZcimp";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}