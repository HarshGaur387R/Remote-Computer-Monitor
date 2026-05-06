const countActiveComputers = (computers: any[]) => {
	let count = 0;
	computers.forEach((value, index) => {
		if (value.active) count = count + 1;
	})
	return count
}


export {
	countActiveComputers
}
