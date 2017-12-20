# GitHub Access Control

Define GitHub teams and team memberships in yaml, render scripts and declarative configurations using templates.

# Why

Initial approach was to define teams as Terraform module which accepts a map or a list to generate team_membership resources.

As team_memberships change, terraform state would cascade when items of the list are removed.

This script is to generate TF Configs without lists, and to simplify managing these configurations by keeping the team
definition in a simple yaml data structure.
