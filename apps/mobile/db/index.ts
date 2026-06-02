import * as SQLite from 'expo-sqlite';
import databaseConstant from '@/constants/database';
const tableName = databaseConstant.tableName;

// Open or create a database file named "computers.db"
const db = SQLite.openDatabaseSync('rcmcomputers.db');

// Create a table on initialization
function initDB() {
	return db.execAsync(`
    CREATE TABLE IF NOT EXISTS computers (
      	id INTEGER PRIMARY KEY AUTOINCREMENT,
      	name VARCHAR(100),
	LANIP VARCHAR(50) NOT NULL,
      	port INTEGER NOT NULL,
	authtoken NOT NULL UNIQUE,
	active BOOLEAN NOT NULL DEFAULT 0
    );
  `);
}

interface ComputerType {
	id: number
	name: string | null
	LANIP: string
	port: number
	authtoken: string
	active: boolean
}

const getTable = async () => {
	try {
		const allrows = await db.getAllAsync(`SELECT * FROM ${tableName}`);
		return allrows as ComputerType[]

		//throw new Error("Testing error")
	} catch (error) {
		throw error
	}
}

const getRowById = async ({ rowId }: { rowId: string }) => {

}

const getRowByIndex = async ({ index }: { index: number }) => {

}

const updateRowById = async (id: number, column: "name" | "LANIP" | "port" | "active", new_value: string | number | boolean) => {
	try {
		await db.runAsync(
			`UPDATE computers SET ${column} = ? WHERE id = ?;`,
			[new_value, id]
		)
	} catch (error) {
		throw error
	}
}

const insertRow = async ({ name, LANIP, authtoken, port, active = false }: Omit<ComputerType, "id">) => {
	try {
		await db.runAsync(
			`INSERT OR IGNORE INTO computers (name, LANIP, authtoken, port, active)
     VALUES (?, ?, ?, ?, ?)`,
			[name, LANIP, authtoken, port, active]
		);

	} catch (error) {
		throw error
	}

}

const deleteRowById = async (id: number) => {
	try {
		await db.runAsync(
			`DELETE FROM computers WHERE id = ?;`,
			[id]
		)
	} catch (error) {
		throw error
	}
}

export {
	initDB,
	getTable,
	ComputerType,
	insertRow,
	updateRowById,
	deleteRowById
}
