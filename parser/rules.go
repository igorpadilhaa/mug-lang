package parser

type RuleType string

const (
	RULE_FUNCALL             RuleType = "RuleFunctionCall"
	RULE_ARGUMENT_LIST       RuleType = "RuleArgumentList"
	RULE_VALUE               RuleType = "RuleValue"
	RULE_VARIABLE            RuleType = "RuleVariable"
	RULE_VARIABLE_ASSIGNMENT RuleType = "RuleVariableAssignment"
)
