package def

const ProjectCardConnectionFragments = `
	fragment projectCardConnection on ProjectCardConnection {
		edges {
			cursor
			node {
				content {
					...issue
				}
				id
				note
				isArchived
			}
		}
	}
`

