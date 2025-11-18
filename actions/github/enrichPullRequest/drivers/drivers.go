package drivers

const BranchName = "branch-name"
const Jira = "jira"
const General = "general"

func Validate(driver string) bool {
    validDrivers := []string{
        BranchName,
        Jira,
        General,
    }

    for _, validDriver := range validDrivers {
        if driver == validDriver {
            return true
        }
    }

    return false
}
