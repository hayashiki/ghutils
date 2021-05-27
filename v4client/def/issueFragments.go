package def

const IssueFragments = `
	fragment issue on Issue {
		title
		number
		url
		body
		labels(first: 10) {
			edges {
				cursor
				node {
					color
					name
				}
			}
		}
		assignees(first: 10) {
			edges {
				cursor
				node {
					id
					name
					url
					avatarUrl(size: 10)
					login
				}
			}
		}
	}
`
