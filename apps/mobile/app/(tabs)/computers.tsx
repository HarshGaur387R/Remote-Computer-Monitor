import { useEffect, useState } from 'react';
import { FlatList, View, Pressable, StyleSheet, Dimensions, ActivityIndicator, TextInput } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { PcDetailedCard } from '@/components/pc-detailed-card';
import { getTable, ComputerType, updateRowById, deleteRowById } from '@/db/index';
import { COLORS } from '@rcm/shared';
import { router } from 'expo-router';
import MaterialIcons from '@expo/vector-icons/MaterialIcons';

const { width } = Dimensions.get('window');
const cardSize = (width - 48) / 2;
const USE_MOCK_DATA = false;

const MOCK_COMPUTERS: ComputerType[] = [
	{ id: 1, name: 'Office Desktop', LANIP: '192.168.1.100', port: 5900, authtoken: 'token_office_desktop_001', active: true },
	{ id: 2, name: 'Development Laptop', LANIP: '192.168.1.105', port: 5901, authtoken: 'token_dev_laptop_002', active: true },
	{ id: 3, name: 'Server Room PC', LANIP: '192.168.1.50', port: 5902, authtoken: 'token_server_room_003', active: false },
	{ id: 4, name: 'Testing Machine', LANIP: '192.168.1.110', port: 5903, authtoken: 'token_testing_machine_004', active: true },
	{ id: 5, name: 'Workstation Pro', LANIP: '192.168.1.115', port: 5904, authtoken: 'token_workstation_pro_005', active: true },
];

// ─── ComputerCard ────────────────────────────────────────────────────────────

interface ComputerCardProps {
	computer: ComputerType;
	onPress: () => void;
	onRename: (id: number, newName: string) => Promise<void>;
}

function ComputerCard({ computer, onPress, onRename }: ComputerCardProps) {
	const [isEditing, setIsEditing] = useState(false);
	const [editedName, setEditedName] = useState(computer.name || 'Unnamed Device');
	const [isSaving, setIsSaving] = useState(false);

	const handleSave = async () => {
		const trimmed = editedName.trim();
		if (!trimmed) return;
		try {
			setIsSaving(true);
			await onRename(computer.id, trimmed);
			setIsEditing(false);
		} catch {
			setEditedName(computer.name || 'Unnamed Device');
		} finally {
			setIsSaving(false);
		}
	};

	const handleCancel = () => {
		setEditedName(computer.name || 'Unnamed Device');
		setIsEditing(false);
	};

	return (
		<Pressable
			style={[styles.card, { width: cardSize, minHeight: cardSize }]}
			onPress={isEditing ? undefined : onPress}
		>
			<View style={styles.cardContent}>

				{/* ── Name row ── */}
				{isEditing ? (
					<View style={styles.cardEditRow}>
						<TextInput
							style={styles.cardNameInput}
							value={editedName}
							onChangeText={setEditedName}
							autoFocus
							selectTextOnFocus
							editable={!isSaving}
							placeholderTextColor="#555"
							onSubmitEditing={handleSave}
							returnKeyType="done"
						/>
						<Pressable
							style={styles.cardIconBtn}
							onPress={handleSave}
							disabled={isSaving}
							hitSlop={{ top: 6, bottom: 6, left: 6, right: 6 }}
						>
							<MaterialIcons
								name={isSaving ? 'hourglass-bottom' : 'check'}
								size={14}
								color={isSaving ? '#555' : '#4ade80'}
							/>
						</Pressable>
						<Pressable
							style={styles.cardIconBtn}
							onPress={handleCancel}
							disabled={isSaving}
							hitSlop={{ top: 6, bottom: 6, left: 6, right: 6 }}
						>
							<MaterialIcons name="close" size={14} color="#888" />
						</Pressable>
					</View>
				) : (
					<View style={styles.cardNameRow}>
						<ThemedText style={styles.computerName} numberOfLines={1}>
							{computer.name || 'Unnamed Device'}
						</ThemedText>
						<Pressable
							onPress={(e) => { e.stopPropagation?.(); setIsEditing(true); }}
							hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
							style={styles.cardPencilBtn}
						>
							<MaterialIcons name="edit" size={13} color={COLORS.primary} />
						</Pressable>
					</View>
				)}

				<ThemedText style={computer.active ? styles.runningState : styles.stopState}>
					{computer.active ? 'running' : 'stopped'}
				</ThemedText>
				<View style={styles.cardDivider} />
				<ThemedText style={styles.cardLabel}>IP</ThemedText>
				<ThemedText style={styles.cardValue} numberOfLines={1}>{computer.LANIP}</ThemedText>
				<ThemedText style={styles.cardLabel}>Port</ThemedText>
				<ThemedText style={styles.cardValue}>{computer.port}</ThemedText>
			</View>
		</Pressable>
	);
}

// ─── Screen ──────────────────────────────────────────────────────────────────

export default function Computers() {
	const [computers, setComputers] = useState<ComputerType[]>([]);
	const [loading, setLoading] = useState(true);
	const [showPcDetails, setShowPcDetails] = useState<{ show: boolean; id: number | null }>({
		show: false,
		id: null,
	});

	useEffect(() => { loadComputers(); }, []);

	const loadComputers = async () => {
		try {
			setLoading(true);
			await new Promise(resolve => setTimeout(resolve, 500));
			const data = USE_MOCK_DATA ? MOCK_COMPUTERS : await getTable();
			setComputers(data);
		} catch (error) {
			console.error('Error loading computers:', error);
		} finally {
			setLoading(false);
		}
	};

	const handleRename = async (id: number, newName: string) => {
		await updateRowById(id, "name", newName);
		setComputers(prev =>
			prev.map(c => c.id === id ? { ...c, name: newName } : c)
		);
	};

	const handleDelete = async (id: number) => {
		await deleteRowById(id);
		setComputers(prev => prev.filter(c => c.id !== id));
	};

	return (
		<SafeAreaView style={styles.safeContainer}>
			<ThemedView style={styles.container}>

				{/* ── Header row ── */}
				<View style={styles.header}>
					<ThemedText style={styles.title}>Computers</ThemedText>
					<Pressable
						onPress={loadComputers}
						disabled={loading}
						style={({ pressed }) => [
							styles.reloadBtn,
							pressed && styles.reloadBtnPressed,
						]}
					>
						{loading ? (
							<ActivityIndicator size="small" color={COLORS.primary} />
						) : (
							<ThemedText style={styles.reloadIcon}>↻</ThemedText>
						)}
					</Pressable>
				</View>

				{loading ? (
					<View style={styles.centerContent}>
						<ActivityIndicator size="large" color={COLORS.primary} />
					</View>
				) : computers.length === 0 ? (
					<View style={styles.centerContent}>
						<ThemedText style={styles.emptyText}>No computers added yet</ThemedText>
						<ThemedText style={styles.emptySubText}>Tap the + button to add a computer</ThemedText>
					</View>
				) : (
					<FlatList
						data={computers}
						keyExtractor={(item) => item.id.toString()}
						numColumns={2}
						columnWrapperStyle={styles.columnWrapper}
						contentContainerStyle={styles.listContent}
						renderItem={({ item }) => (
							<ComputerCard
								computer={item}
								onPress={() => setShowPcDetails({ show: true, id: item.id })}
								onRename={handleRename}
							/>
						)}
						scrollEnabled
					/>
				)}

				{/* FAB */}
				<Pressable
					style={[styles.fab, { backgroundColor: COLORS.primary }]}
					onPress={() => { router.push('/qrscanner'); }}
				>
					<ThemedText style={styles.fabText}>+</ThemedText>
				</Pressable>

				{/* Detail modal — at ThemedView level, never inside FAB */}
				{showPcDetails.show && showPcDetails.id !== null && (
					<PcDetailedCard
						computerId={showPcDetails.id}
						onClose={() => setShowPcDetails({ show: false, id: null })}
						onUpdateComputerName={handleRename}
						onDelete={handleDelete}
					/>
				)}
			</ThemedView>
		</SafeAreaView>
	);
}

// ─── Styles ──────────────────────────────────────────────────────────────────

const styles = StyleSheet.create({
	safeContainer: {
		flex: 1,
		backgroundColor: '#0C0C0C',
	},
	container: {
		flex: 1,
		paddingHorizontal: 10,
		paddingTop: 8,
		position: 'relative',
		backgroundColor: 'transparent',
	},

	// ── NEW: header row ──
	header: {
		flexDirection: 'row',
		alignItems: 'center',
		justifyContent: 'space-between',
		marginBottom: 12,
		marginLeft: 6,
		marginRight: 6,
	},
	title: {
		fontSize: 20,
		fontWeight: '600',
		color: '#888',
	},
	reloadBtn: {
		width: 34,
		height: 34,
		borderRadius: 8,
		backgroundColor: 'rgba(15, 40, 84, 0.35)',
		borderWidth: 0.5,
		borderColor: 'rgba(15, 40, 84, 0.6)',
		justifyContent: 'center',
		alignItems: 'center',
	},
	reloadBtnPressed: {
		opacity: 0.6,
	},
	reloadIcon: {
		fontSize: 20,
		color: '#888',
		lineHeight: 22,
	},

	centerContent: {
		flex: 1,
		justifyContent: 'center',
		alignItems: 'center',
		paddingHorizontal: 32,
	},
	emptyText: {
		fontSize: 18,
		fontWeight: '600',
		marginBottom: 8,
		textAlign: 'center',
		color: 'white',
	},
	emptySubText: {
		fontSize: 14,
		opacity: 0.6,
		textAlign: 'center',
		color: '#888',
	},
	columnWrapper: {
		justifyContent: 'space-between',
		marginBottom: 16,
		paddingHorizontal: 6,
	},
	listContent: {
		paddingBottom: 100,
	},

	// Card styles unchanged…
	card: {
		backgroundColor: 'rgba(15, 40, 84, 0.3)',
		borderRadius: 12,
		borderWidth: 0.5,
		borderColor: 'rgba(15, 40, 84, 0.5)',
		padding: 12,
		justifyContent: 'flex-start',
		shadowColor: '#000',
		shadowOffset: { width: 0, height: 2 },
		shadowOpacity: 0.1,
		shadowRadius: 4,
		elevation: 3,
	},
	cardContent: { flex: 1, justifyContent: 'flex-start' },
	cardNameRow: { flexDirection: 'row', alignItems: 'center', marginBottom: 6 },
	computerName: { fontSize: 14, fontWeight: '600', color: 'white', flex: 1 },
	cardPencilBtn: { marginLeft: 4, padding: 2 },
	cardEditRow: { flexDirection: 'row', alignItems: 'center', marginBottom: 6, gap: 4 },
	cardNameInput: {
		flex: 1,
		fontSize: 13,
		fontWeight: '600',
		color: 'white',
		backgroundColor: 'rgba(15, 40, 84, 0.4)',
		borderWidth: 1,
		borderColor: COLORS.primary,
		borderRadius: 6,
		paddingHorizontal: 7,
		paddingVertical: 3,
	},
	cardIconBtn: { padding: 4, borderRadius: 4, backgroundColor: 'rgba(255,255,255,0.05)' },
	runningState: { color: '#4ade80', fontSize: 11, fontWeight: '500' },
	stopState: { color: '#ef4444', fontSize: 11, fontWeight: '500' },
	cardDivider: { height: 0.5, backgroundColor: 'rgba(15, 40, 84, 0.5)', marginVertical: 8 },
	cardLabel: { fontSize: 11, color: '#888', marginBottom: 2, fontWeight: '500' },
	cardValue: { fontSize: 12, fontWeight: '500', color: '#ccc', marginBottom: 4 },

	// FAB
	fab: {
		position: 'absolute',
		bottom: 24,
		right: 24,
		width: 56,
		height: 56,
		borderRadius: 28,
		justifyContent: 'center',
		alignItems: 'center',
		shadowColor: '#000',
		shadowOffset: { width: 0, height: 4 },
		shadowOpacity: 0.3,
		shadowRadius: 8,
		elevation: 8,
	},
	fabText: { fontSize: 28, fontWeight: '300', color: '#fff' },
});
