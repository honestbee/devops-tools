#!/usr/bin/env python

import hcl

print_var_format="|`{parameter:<25}`|`{default:<20}`|{description:<80}|{required:<8}|"
print_out_format="|`{parameter:<25}`|{description:<80}|"

var_header="|{name:^27}|{default:^22}|{description:^80}|{required:^8}|\n|{a:-<27}|{a:-<22}|{a:-<80}|{a:-<8}|".format(
        name="Name",
        default="Default",
        description="Description",
        required="Required",
        a=":"
)

out_header="|{name:<27}|{description:<80}|\n|{a:-<27}|{a:-<80}|".format(
        name="Name",
        description="Description",
        a=":"
)

with open('variables.tf', 'r') as fp:
    d = hcl.load(fp)
    var = d.get("variable")
    if var != None:
        print("## Variables\n")
        print(var_header)
        for k,v in sorted(var.items(),key=lambda x:x[0]):
            required = v.get('default') == None
            s = print_var_format.format(
                    parameter=k,
                    default=v.get('default'),
                    description=v.get('description'),
                    required= "Yes" if required else "No"
                )
            print(s)
        print("\n")
    out = d.get("output")
    if out != None:
        print("## Outputs\n")
        print(out_header)
        for k,v in sorted(out.items(),key=lambda x:x[0]):
            s = print_out_format.format(
                    parameter=k,
                    description=v.get('description')
                )
            print(s)
        print("\n")

