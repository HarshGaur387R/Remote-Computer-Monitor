import { StyleSheet } from 'react-native';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { SafeAreaView } from 'react-native-safe-area-context';
import AnimatedLoader from '@/components/animated-loader';
import * as Network from 'expo-network';
import { useNetworkState } from 'expo-network';
import { APP_NAME, COLORS } from "@rcm/shared"
import { Loading_screen } from '@/components/LoadingScreen';
import { NetworkErrorScreen } from '@/components/NetworkErrorScreen';
import { ReactNode, useEffect, useState } from 'react';
import { getTable } from '@/db';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { TouchableOpacity } from 'react-native';
import { countActiveComputers } from '@/utils';
import { ScrollView } from 'react-native';
import { ErrorScreen } from '@/components/ErrorScreen';
import { router } from 'expo-router';

function FrameContainer({ child }: { child: ReactNode }) {
	return (
		<SafeAreaView style={styles.container}>
			<ThemedView style={styles.mainScreen}>
				<ThemedText style={styles.title}>
					{APP_NAME}
				</ThemedText>
				{
					child
				}

			</ThemedView>
		</SafeAreaView>
	)
}

type Computers = {
	id: number;
	name: string;
	ip: string;
	active: boolean;
}

type DisplayComputerInformationType = {
	connectedComputers: number;
	activeComputers: number;
	computers: { id: number; name: string; ip: string; active: boolean }[];
}

function DisplayComputerInformation({ connectedComputers, activeComputers, computers }: DisplayComputerInformationType) {
	return (
		<ThemedView style={styles.displayComputerInformationContainer}>

			{/* Summary cards */}
			<ThemedView style={styles.metricsRow}>
				<ThemedView style={styles.metricCard}>
					<ThemedText style={styles.metricLabel}>Connected</ThemedText>
					<ThemedText style={styles.metricValue}>{connectedComputers}</ThemedText>
					<ThemedText style={styles.metricSub}>computers</ThemedText>
				</ThemedView>

				<ThemedView style={styles.metricCard}>
					<ThemedView style={styles.activeDotRow}>
						<ThemedView style={styles.activeDot} />
						<ThemedText style={styles.metricLabel}>Active</ThemedText>
					</ThemedView>
					<ThemedText style={styles.metricValue}>{activeComputers}</ThemedText>
					<ThemedText style={styles.metricSub}>of {connectedComputers} online</ThemedText>
				</ThemedView>
			</ThemedView>

			{/* Device list */}
			<ThemedView style={styles.deviceList}>
				<ThemedText style={styles.deviceListHeader}>Devices</ThemedText>
				<ScrollView
					style={styles.deviceScrollView}
					showsVerticalScrollIndicator={false}
					bounces={false}
				>
					{computers.map((pc, i) => (
						<ThemedView
							key={pc.id}
							style={[styles.deviceRow, i < computers.length - 1 && styles.deviceRowBorder]}
						>
							<ThemedView style={styles.deviceIcon}>
								{/* Replace with your IconSymbol */}
								<IconSymbol name="pc" size={20} color={COLORS.secondery} />
							</ThemedView>
							<ThemedView style={styles.deviceInfo}>
								<ThemedText style={styles.deviceName}>{pc.name}</ThemedText>
								<ThemedText style={styles.deviceIp}>{pc.ip}</ThemedText>
							</ThemedView>
							<ThemedView style={[styles.badge, pc.active ? styles.badgeActive : styles.badgeIdle]}>
								<ThemedText style={[styles.badgeText, pc.active ? styles.badgeTextActive : styles.badgeTextIdle]}>
									{pc.active ? 'Active' : 'Idle'}
								</ThemedText>
							</ThemedView>
						</ThemedView>
					))}
				</ScrollView>
			</ThemedView>

			<ThemedView style={styles.messageContainer}>
				<ThemedText style={styles.messageText}>
					Sometimes each PC can have different LOCAL IP Address because
					they are not static address which is associated to a specific pc.
				</ThemedText>
			</ThemedView>

		</ThemedView>
	);
}

function DiplayAddComputer() {
	return (
		<ThemedView style={styles.displayAddComputerContainer}>
			<IconSymbol name='qrcode' size={130} color={"#1C4D8D"} />
			<ThemedText style={{
				textAlign: 'center',
				fontSize: 20,
				color: 'darkgrey'
			}}>
				{"No connected computers found!"}
			</ThemedText>
			<ThemedText style={{
				textAlign: 'center',
				color: 'grey'
			}}>
				{"Scan client's QR code to connect with agent"}
			</ThemedText>
			<TouchableOpacity
				style={styles.button}
				onPress={() => { router.push('/qrscanner')}}
				activeOpacity={0.7}
			>
				<ThemedText style={styles.buttonText}>{"Scan QR"}</ThemedText>
			</TouchableOpacity>
		</ThemedView>
	)
}


export default function HomeScreen() {
	console.log("HomeScreen renderd!")
	const networkState = (useNetworkState()).type === Network.NetworkStateType.WIFI;
	const [loading, setLoading] = useState<boolean>(true)
	const [error, setError] = useState<string | null>(null)
	const [activeComputers, setActiveComputers] = useState<number>(0)
	const [connectedComputers, setConnectedComputers] = useState<number>(0)
	const [computers, setComputers] = useState<unknown[]>([])

	useEffect(() => {
		async function fetchTableRows() {
			try {
				setError(null)
				setLoading(true);
				const rows = await getTable("computers")

				setComputers(rows as Computers[])
				const activeCount = countActiveComputers(computers);
				const connectedCounts = computers.length;

				setActiveComputers(activeCount);
				setConnectedComputers(connectedCounts);
				setLoading(false)
			} catch (error) {
				const err = error as Error;
				setError(err.message);
			} finally {
				setLoading(false);
			}
		}
		fetchTableRows()
	}, [])

	if (error) {
		return (
			<FrameContainer child={
				<ErrorScreen
					container_style={styles.networkErrorConainer}
					error='Error Detected'
					message={error}
				/>
			} />

		)
	}

	if (!networkState) {
		return (
			<FrameContainer child={
				<NetworkErrorScreen
					container_style={styles.networkErrorConainer}
					error='WIFI is OFF'
					message='WiFi is not connected. Please turn it ON and connected to the router'
				/>
			} />

		)
	} else if (loading) {
		return (<FrameContainer child={
			<Loading_screen
				container_style={styles.loaderContainer}
				animated_loader_child={<AnimatedLoader type='loader1' size={300} />}
				text={"Looking for Computers..."}
			/>}
		/>)
	} else {
		return (
			<FrameContainer child=
				{
					computers.length > 0
						? <DisplayComputerInformation
							connectedComputers={connectedComputers}
							activeComputers={activeComputers}
							computers={computers as Computers[]}

						/>
						: <DiplayAddComputer />
				}
			/>
		)
	}
}

const styles = StyleSheet.create({
	container: {
		padding: 5,
		width: "100%",
		height: "100%",
		flex: 1,
		alignItems: 'center',
		justifyContent: 'center',
		backgroundColor: "#0C0C0C",
		//borderWidth: 2,
		//borderColor: 'red',
	},
	title: {
		fontSize: 28,
		color: 'white',
		textAlign: 'center',
		padding: 8,
		fontWeight: 'semibold',
		lineHeight: 30,
	},
	mainScreen: {
		padding: 10,
		width: "90%",
		height: 600,
		backgroundColor: 'transparent',
		display: 'flex',
		flexDirection: 'column',
		//borderWidth: 2,
		//borderColor: 'blue'
	},
	loaderContainer: {
		height: "80%",
		justifyContent: 'center',
		alignItems: 'center',
		backgroundColor: 'transparent'
	},
	networkErrorConainer: {
		height: "100%",
		justifyContent: 'center',
		alignItems: 'center',
		backgroundColor: 'transparent',
	},
	displayComputerInformationContainer: {
		//borderWidth: 2,
		//borderColor: 'yellow',
		height: "90%",
		justifyContent: "flex-start",
		alignItems: "center",
		backgroundColor: "transparent"
	},
	informationCard: {
		width: "100%",
		borderColor: 'green',
		borderWidth: 2
	},
	displayAddComputerContainer: {
		//borderWidth: 2,
		//borderColor: 'yellow',
		height: "90%",
		justifyContent: "center",
		alignItems: "center",
		backgroundColor: "transparent"
	},
	button: {
		width: 120,
		paddingVertical: 10,
		marginTop: 20,
		backgroundColor: '#1C4D8D',
		borderRadius: 8,
		alignItems: 'center',
	},
	buttonText: {
		color: 'white',
		fontSize: 14,
		fontWeight: '500',
	},
	metricsRow: {
		flexDirection: 'row',
		gap: 12,
		marginBottom: 16,
		backgroundColor: 'transparent',
	},
	metricCard: {
		flex: 1,
		//backgroundColor: '#1A1A1A',
		backgroundColor: 'rgb(15, 40, 84,0.3)',
		borderRadius: 12,
		padding: 16,
		gap: 4,
	},
	activeDotRow: {
		flexDirection: 'row',
		alignItems: 'center',
		gap: 6,
		backgroundColor: 'transparent',
	},
	activeDot: {
		width: 8,
		height: 8,
		borderRadius: 4,
		backgroundColor: '#22c55e',
	},
	metricLabel: { fontSize: 13, color: '#888' },
	metricValue: { fontSize: 32, fontWeight: '500', color: 'white', lineHeight: 38 },
	metricSub: { fontSize: 12, color: '#555' },

	deviceList: {
		width: "100%",
		//backgroundColor: '#111',
		backgroundColor: 'rgb(15, 40, 84,0.3)',
		borderRadius: 12,
		borderWidth: 0.5,
		borderColor: 'rgb(15, 40, 84,0.5)',
		overflow: 'hidden',
	},
	deviceListHeader: {
		fontSize: 13,
		fontWeight: '500',
		color: '#888',
		paddingHorizontal: 16,
		paddingVertical: 10,
		borderBottomWidth: 0.5,
		borderBottomColor: '#2a2a2a',
	},
	deviceRow: {
		flexDirection: 'row',
		alignItems: 'center',
		gap: 12,
		paddingHorizontal: 16,
		paddingVertical: 12,
		backgroundColor: 'transparent',
	},
	deviceRowBorder: { borderBottomWidth: 0.5, borderBottomColor: '#2a2a2a' },
	deviceIcon: {
		width: 36, height: 36, borderRadius: 8,
		backgroundColor: COLORS.primary,
		alignItems: 'center', justifyContent: 'center',
	},
	deviceInfo: { flex: 1, backgroundColor: 'transparent' },
	deviceName: { fontSize: 14, fontWeight: '500', color: 'white' },
	deviceIp: { fontSize: 12, color: '#666' },

	badge: { borderRadius: 99, paddingHorizontal: 10, paddingVertical: 3 },
	badgeActive: { backgroundColor: '#14532d22' },
	badgeIdle: { backgroundColor: '#1e293b' },
	badgeText: { fontSize: 11 },
	badgeTextActive: { color: '#4ade80' },
	badgeTextIdle: { color: '#94a3b8' },
	messageContainer: {
		backgroundColor: 'transparent',
	},
	messageText: {
		fontSize: 12,
		color: '#666',
		lineHeight: 16,
		padding: 10,
		paddingTop: 15,
	},
	deviceScrollView: {
		maxHeight: 280, // fits ~4-5 rows before scrolling kicks in, adjust to taste
	},

});
