# Spacelift Plan Approval Policy
# Controls when runs require manual approval vs auto-approve

package spacelift

# Default: require approval
default approve = false

# Auto-approve runs on non-foundation stacks with no destroys
approve {
    not is_foundation_stack
    not has_destroys
    not is_administrative
}

# Auto-approve task runs (like terraform destroy commands)
approve {
    input.run.type == "TASK"
}

# Helper: Check if this is a foundation stack
is_foundation_stack {
    input.stack.labels[_] == "foundation"
}

# Helper: Check if this is an administrative stack
is_administrative {
    input.stack.administrative == true
}

# Helper: Check if the plan includes resource deletions
has_destroys {
    input.run.changes.deleted > 0
}

# Warn on foundation changes
warn["Foundation stack changes should be reviewed carefully"] {
    is_foundation_stack
}

# Warn on administrative stack changes
warn["Administrative stack changes affect all other stacks"] {
    is_administrative
}

# Warn on large changes
warn[sprintf("Large change: %d resources will be modified", [total])] {
    total := input.run.changes.added + input.run.changes.changed + input.run.changes.deleted
    total > 10
}

# Deny changes to certain protected resources (example)
deny["Cannot modify production databases without explicit approval"] {
    input.run.changes.changed > 0
    resource := input.terraform.resource_changes[_]
    resource.type == "azurerm_postgresql_flexible_server"
    resource.change.actions[_] == "delete"
}

# Sample notification for cost tracking
sample {
    true
}
