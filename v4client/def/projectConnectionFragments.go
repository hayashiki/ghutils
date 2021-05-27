package def

const ProjectConnectionFragments = `
	fragment projectConnection on ProjectConnection {
		edges {
			cursor
			node {
				id
				url
				name
				columns(first: 10) {
					edges {
						cursor
						node {
							name
							id
							cards(first: 100) {
								...projectCardConnection
							}
						}
					}
				}
			}
		}
	}
`
