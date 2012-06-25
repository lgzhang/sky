#include <stdio.h>
#include <stdlib.h>

#include <eql/ast.h>
#include <eql/parser.h>
#include <eql/compiler.h>

#include "../minunit.h"


//==============================================================================
//
// Globals
//
//==============================================================================

struct tagbstring foo = bsStatic("foo");


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// AST
//--------------------------------------

int test_eql_ast_float_literal_create() {
    eql_ast_node *node;
    eql_ast_float_literal_create(100.293, &node);
    mu_assert(node->type == EQL_AST_TYPE_FLOAT_LITERAL, "");
    mu_assert(node->float_literal.value == 100.293, "");
    mu_assert_eql_node_dump(node, "<float-literal value='100.29300'>\n");
    eql_ast_node_free(node);
    return 0;
}


//--------------------------------------
// Parser
//--------------------------------------

int test_eql_float_literal_parse() {
    eql_ast_node *module = NULL;
    bstring text = bfromcstr("10.421;");
    eql_parse(NULL, text, &module);
    eql_ast_node *node = module->module.main_function->function.body->block.exprs[0];
    mu_assert(node->type == EQL_AST_TYPE_FLOAT_LITERAL, "");
    mu_assert(node->float_literal.value == 10.421, "");
    eql_ast_node_free(module);
    bdestroy(text);
    return 0;
}


//--------------------------------------
// Type
//--------------------------------------

int test_eql_float_literal_get_type() {
    eql_ast_node *node;
    bstring type;
    eql_module *module = eql_module_create(&foo, NULL);
    eql_ast_float_literal_create(100.1, &node);
    eql_ast_node_get_type(node, module, &type);
    mu_assert(biseqcstr(type, "Float"), "");
    eql_ast_node_free(node);
    eql_module_free(module);
    bdestroy(type);
    return 0;
}


//--------------------------------------
// Compile
//--------------------------------------

int test_eql_compile_float_literal() {
    mu_assert_eql_compile("return 10.25;", "tests/fixtures/eql/ir/float_literal.ll");
    return 0;
}

//--------------------------------------
// Execute  
//--------------------------------------

int test_eql_execute_float_literal() {
    mu_assert_eql_execute_float("return 10.25;", 10.25);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_eql_ast_float_literal_create);
    mu_run_test(test_eql_float_literal_parse);
    mu_run_test(test_eql_float_literal_get_type);
    mu_run_test(test_eql_compile_float_literal);
    mu_run_test(test_eql_execute_float_literal);
    return 0;
}

RUN_TESTS()