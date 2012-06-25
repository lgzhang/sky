#ifndef _eql_ast_function_h
#define _eql_ast_function_h

#include "../../bstring.h"

//==============================================================================
//
// Definitions
//
//==============================================================================

// Represents a function in the AST.
typedef struct {
    bstring name;
    bstring return_type;
    struct eql_ast_node **args;
    unsigned int arg_count;
    struct eql_ast_node *body;
} eql_ast_function;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

int eql_ast_function_create(bstring name, bstring return_type,
    struct eql_ast_node **args, unsigned int arg_count,
    struct eql_ast_node *body, struct eql_ast_node **ret);

void eql_ast_function_free(struct eql_ast_node *node);


//--------------------------------------
// Codegen
//--------------------------------------

int eql_ast_function_codegen(eql_ast_node *node, eql_module *module,
    LLVMValueRef *value);

int eql_ast_function_codegen_args(eql_ast_node *node, eql_module *module);


//--------------------------------------
// Misc
//--------------------------------------

int eql_ast_function_get_class(eql_ast_node *node, eql_ast_node **class_ast);

int eql_ast_function_generate_return_type(eql_ast_node *node, eql_module *module);

int eql_ast_function_get_var_decl(eql_ast_node *node, bstring name,
    eql_ast_node **var_decl);

//--------------------------------------
// Debugging
//--------------------------------------

int eql_ast_function_dump(eql_ast_node *node, bstring ret);

#endif