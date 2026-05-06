import * as SQLite from 'expo-sqlite';
// Open or create a database file named "computers.db"
const db = SQLite.openDatabaseSync('rcmcomputers.db');
// Create a table

function initDB() {
	db.execAsync(`
    CREATE TABLE IF NOT EXISTS computers (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT,
      ip TEXT
    );
  `);
}

const getTable = async (tableName: string) => {
	try {
		const allrows = await db.getAllAsync(`SELECT * FROM ${tableName}`);
		return allrows

		//throw new Error("Testing error")


		//return [
		//	{ id: 1, name: "PC1", ip: "192.168.1.10", active: false },
		//	{ id: 2, name: "PC2", ip: "192.168.1.11", active: true },
		//	{ id: 3, name: "PC3", ip: "192.168.1.12", active: true },
		//	{ id: 4, name: "PC4", ip: "192.168.1.12", active: true },
		//	{ id: 5, name: "PC5", ip: "192.168.1.12", active: true },
		//	{ id: 6, name: "PC6", ip: "192.168.1.12", active: true },
		//	{ id: 7, name: "PC7", ip: "192.168.1.12", active: true },
		//]

	} catch (error) {
		throw error
	}
}

const getRowById = async ({ tableName, rowId }: { tableName: string, rowId: string }) => {

}

const getRowByindex = async ({ tableName, index }: { tableName: string, index: number }) => {

}

export {
	initDB,
	getTable
}
