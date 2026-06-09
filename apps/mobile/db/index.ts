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
	active BOOLEAN NOT NULL DEFAULT 0,
	cpu_usage_percent FLOAT,
	total_vram_mb INTEGER,
	available_vram_mb INTEGER,
	vram_used_percent FLOAT,
	total_storage_gb INTEGER,
	available_storage_gb INTEGER,
	storage_used_percent FLOAT
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
	cpu_usage_percent?: number | null
	total_vram_mb?: number | null
	available_vram_mb?: number | null
	vram_used_percent?: number | null
	total_storage_gb?: number | null
	available_storage_gb?: number | null
	storage_used_percent?: number | null
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

const updateRowById = async (id: number, column: "name" | "LANIP" | "port" | "active" | "cpu_usage_percent" | "total_vram_mb" | "available_vram_mb" | "vram_used_percent" | "total_storage_gb" | "available_storage_gb" | "storage_used_percent", new_value: string | number | boolean | null) => {
	try {
		await db.runAsync(
			`UPDATE computers SET ${column} = ? WHERE id = ?;`,
			[new_value, id]
		)
	} catch (error) {
		throw error
	}
}

const insertRow = async ({ name, LANIP, authtoken, port, active = false, cpu_usage_percent, total_vram_mb, available_vram_mb, vram_used_percent, total_storage_gb, available_storage_gb, storage_used_percent }: Omit<ComputerType, "id"> & { cpu_usage_percent?: number | null; total_vram_mb?: number | null; available_vram_mb?: number | null; vram_used_percent?: number | null; total_storage_gb?: number | null; available_storage_gb?: number | null; storage_used_percent?: number | null }) => {
	try {
		await db.runAsync(
			`INSERT OR IGNORE INTO computers (name, LANIP, authtoken, port, active, cpu_usage_percent, total_vram_mb, available_vram_mb, vram_used_percent, total_storage_gb, available_storage_gb, storage_used_percent)
     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			[name, LANIP, authtoken, port, active, cpu_usage_percent ?? null, total_vram_mb ?? null, available_vram_mb ?? null, vram_used_percent ?? null, total_storage_gb ?? null, available_storage_gb ?? null, storage_used_percent ?? null]
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

const updateFieldsById = async (id: number, updates: Partial<Omit<ComputerType, 'id' | 'authtoken'>>) => {
	try {
		// Filter out id and authtoken if they somehow exist in updates
		const allowedUpdates = Object.entries(updates).filter(
			([key]) => key !== 'id' && key !== 'authtoken'
		);

		if (allowedUpdates.length === 0) {
			return; // No valid fields to update
		}

		// Build the SET clause dynamically
		const setClause = allowedUpdates.map(([key]) => `${key} = ?`).join(', ');
		const values = allowedUpdates.map(([, value]) => value);

		await db.runAsync(
			`UPDATE computers SET ${setClause} WHERE id = ?;`,
			[...values, id]
		);
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
	updateFieldsById,
	deleteRowById
}
