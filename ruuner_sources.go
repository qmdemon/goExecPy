package goExecPy

const (
	Python_source string = `
import json


#{source}


def exec_python_custom_functions(custom_function,custom_args=''):
    global_namespace = globals()
    if custom_function in global_namespace:
        custom_functions_string = f"{custom_function}({','.join(json.dumps(i) for i in custom_args)})"
    else:
        custom_functions_string = custom_function
    return eval(custom_functions_string)

try:
    custom_functions_result = exec_python_custom_functions(#{encoded_source})
    print('')
    if ( custom_functions_result == None):
        print('["ok"]')
    else:
        try :
            print(json.dumps(['ok', custom_functions_result]))
        except Exception as e:
            # print('["err"]')
            print(json.dumps(['err', '' + str(e)]))
except Exception as e:
    print(json.dumps(['err', '' + str(e)]))

`
)
