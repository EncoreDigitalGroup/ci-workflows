package drivers

const BranchName = "branch-name"
const Jira = "jira"

func Validate(driver string) bool {
    validDrivers := []string{
        BranchName,
        Jira,
    }

    for _, validDriver := range validDrivers {
        if driver == validDriver {
            return true
        }
    }

    return false
}
