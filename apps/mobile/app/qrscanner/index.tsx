import React, { useState, useEffect, useRef } from "react";
import { CameraView, Camera, BarcodeScanningResult } from "expo-camera";
import {
	Text,
	View,
	StyleSheet,
	TouchableOpacity,
	Animated,
	Easing,
} from "react-native";
import { insertRow } from "@/db";

interface ScanDataType {
	lan_ip: string;
	auth_token: string;
	port: number;
}

interface SystemInfoType {
	cpu_usage_percent: number;
	total_vram_mb: number;
	available_vram_mb: number;
	vram_used_percent: number;
	total_storage_gb: number;
	available_storage_gb: number;
	storage_used_percent: number;
}

type ScanState =
	| { status: "idle" }
	| { status: "processing" }
	| { status: "success"; ip: string }
	| { status: "error"; message: string };

// ─── Loading overlay (shown while processing a scanned QR) ────────

function ProcessingOverlay() {
	const dot1 = useRef(new Animated.Value(0)).current;
	const dot2 = useRef(new Animated.Value(0)).current;
	const dot3 = useRef(new Animated.Value(0)).current;

	useEffect(() => {
		const animateDot = (dot: Animated.Value, delay: number) =>
			Animated.loop(
				Animated.sequence([
					Animated.delay(delay),
					Animated.timing(dot, {
						toValue: 1,
						duration: 400,
						easing: Easing.inOut(Easing.ease),
						useNativeDriver: true,
					}),
					Animated.timing(dot, {
						toValue: 0,
						duration: 400,
						easing: Easing.inOut(Easing.ease),
						useNativeDriver: true,
					}),
				])
			);

		const a1 = animateDot(dot1, 0);
		const a2 = animateDot(dot2, 200);
		const a3 = animateDot(dot3, 400);
		a1.start(); a2.start(); a3.start();
		return () => { a1.stop(); a2.stop(); a3.stop(); };
	}, []);

	const dotStyle = (anim: Animated.Value) => ({
		opacity: anim.interpolate({ inputRange: [0, 1], outputRange: [0.3, 1] }),
		transform: [{
			translateY: anim.interpolate({ inputRange: [0, 1], outputRange: [0, -6] }),
		}],
	});

	return (
		<View style={overlayStyles.backdrop}>
			<View style={overlayStyles.card}>
				<Text style={overlayStyles.title}>Connecting to agent</Text>
				<Text style={overlayStyles.subtitle}>Verifying device on local network…</Text>
				<View style={overlayStyles.dotsRow}>
					<Animated.View style={[overlayStyles.dot, dotStyle(dot1)]} />
					<Animated.View style={[overlayStyles.dot, dotStyle(dot2)]} />
					<Animated.View style={[overlayStyles.dot, dotStyle(dot3)]} />
				</View>
			</View>
		</View>
	);
}

// ─── Error overlay (shown when scan processing fails) ─────────────

interface ErrorOverlayProps {
	message: string;
	onRetry: () => void;
}

function ErrorOverlay({ message, onRetry }: ErrorOverlayProps) {
	return (
		<View style={overlayStyles.backdrop}>
			<View style={overlayStyles.card}>
				<Text style={overlayStyles.errorIcon}>⚠</Text>
				<Text style={overlayStyles.title}>Scan Failed</Text>
				<View style={overlayStyles.errorBox}>
					<Text style={overlayStyles.errorMessage} numberOfLines={4}>
						{message}
					</Text>
				</View>
				<TouchableOpacity style={overlayStyles.retryButton} onPress={onRetry} activeOpacity={0.75}>
					<Text style={overlayStyles.retryText}>Scan Again</Text>
				</TouchableOpacity>
			</View>
		</View>
	);
}

// ─── Permission screens ───────────────────────────────────────────

function PermissionRequestingScreen() {
	return (
		<View style={permStyles.container}>
			<Text style={permStyles.icon}>📷</Text>
			<Text style={permStyles.title}>Camera Access Needed</Text>
			<Text style={permStyles.subtitle}>Requesting permission to use the camera…</Text>
		</View>
	);
}

function PermissionDeniedScreen() {
	return (
		<View style={permStyles.container}>
			<Text style={permStyles.icon}>🚫</Text>
			<Text style={permStyles.title}>Camera Access Denied</Text>
			<Text style={permStyles.subtitle}>
				Enable camera permission in your device Settings to scan QR codes.
			</Text>
		</View>
	);
}

// ─── Main component ───────────────────────────────────────────────

export default function QRScannerScreen() {
	const [hasPermission, setHasPermission] = useState<boolean | null>(null);
	const [scanState, setScanState] = useState<ScanState>({ status: "idle" });

	async function fetchSystemInfo(jsonData: ScanDataType): Promise<SystemInfoType> {
		const controller = new AbortController();
		const timeout = setTimeout(() => controller.abort(), 5000);

		try {
			const response = await fetch(
				`http://${jsonData.lan_ip}:${jsonData.port}/system-info?auth=${jsonData.auth_token}`,
				{
					signal: controller.signal,
				}
			);
			if (!response.ok)
				throw new Error(`Agent responded with status ${response.status}`);

			const jsonResponse = await response.json();
			console.log(jsonResponse);

			// Validate required fields
			const systemInfo: SystemInfoType = {
				cpu_usage_percent: typeof jsonResponse.cpu_usage_percent === 'number' ? jsonResponse.cpu_usage_percent : 0,
				total_vram_mb: typeof jsonResponse.total_vram_mb === 'number' ? jsonResponse.total_vram_mb : 0,
				available_vram_mb: typeof jsonResponse.available_vram_mb === 'number' ? jsonResponse.available_vram_mb : 0,
				vram_used_percent: typeof jsonResponse.vram_used_percent === 'number' ? jsonResponse.vram_used_percent : 0,
				total_storage_gb: typeof jsonResponse.total_storage_gb === 'number' ? jsonResponse.total_storage_gb : 0,
				available_storage_gb: typeof jsonResponse.available_storage_gb === 'number' ? jsonResponse.available_storage_gb : 0,
				storage_used_percent: typeof jsonResponse.storage_used_percent === 'number' ? jsonResponse.storage_used_percent : 0,
			};

			return systemInfo;
		} finally {
			clearTimeout(timeout);
		}
	}

	useEffect(() => {
		const getCameraPermissions = async () => {
			const { status } = await Camera.requestCameraPermissionsAsync();
			setHasPermission(status === "granted");
		};
		getCameraPermissions();
	}, []);

	const handleBarcodeScanned = async ({ data }: BarcodeScanningResult) => {
		// Guard immediately to prevent duplicate scans during async work
		setScanState({ status: "processing" });

		try {
			const jsonData = JSON.parse(data) as ScanDataType;

			if (!jsonData.lan_ip || !jsonData.auth_token || typeof jsonData.port !== "number") {
				throw new Error("QR code is missing required fields (lan_ip, auth_token, port)");
			}

			const systemInfo = await fetchSystemInfo(jsonData);
			console.log(systemInfo);

			await insertRow({
				name: null,
				LANIP: jsonData.lan_ip,
				authtoken: jsonData.auth_token,
				port: jsonData.port,
				active: true,
				cpu_usage_percent: systemInfo.cpu_usage_percent,
				total_vram_mb: systemInfo.total_vram_mb,
				available_vram_mb: systemInfo.available_vram_mb,
				vram_used_percent: systemInfo.vram_used_percent,
				total_storage_gb: systemInfo.total_storage_gb,
				available_storage_gb: systemInfo.available_storage_gb,
				storage_used_percent: systemInfo.storage_used_percent,
			});

			setScanState({ status: "success", ip: jsonData.lan_ip });
		} catch (error) {
			const err = error as Error;
			setScanState({ status: "error", message: err.message });
			console.error(error);
		}
	};

	const resetScan = () => setScanState({ status: "idle" });

	if (hasPermission === null) return <PermissionRequestingScreen />;
	if (hasPermission === false) return <PermissionDeniedScreen />;

	const isScanning = scanState.status === "idle";

	return (
		<View style={styles.container}>
			<CameraView
				onBarcodeScanned={isScanning ? handleBarcodeScanned : undefined}
				barcodeScannerSettings={{ barcodeTypes: ["qr", "pdf417"] }}
				style={StyleSheet.absoluteFillObject}
			/>

			{scanState.status === "success" && (
				<View style={overlayStyles.backdrop}>
					<View style={overlayStyles.card}>
						<Text style={overlayStyles.successIcon}>✓</Text>
						<Text style={overlayStyles.title}>Device Added</Text>
						<Text style={overlayStyles.subtitle}>{scanState.ip}</Text>
						<TouchableOpacity style={overlayStyles.retryButton} onPress={resetScan} activeOpacity={0.75}>
							<Text style={overlayStyles.retryText}>Scan Another</Text>
						</TouchableOpacity>
					</View>
				</View>
			)}

			{scanState.status === "processing" && <ProcessingOverlay />}

			{scanState.status === "error" && (
				<ErrorOverlay message={scanState.message} onRetry={resetScan} />
			)}
		</View>
	);
}

// ─── Styles ───────────────────────────────────────────────────────

const styles = StyleSheet.create({
	container: {
		flex: 1,
		flexDirection: "column",
		justifyContent: "center",
	},
});

const overlayStyles = StyleSheet.create({
	backdrop: {
		...StyleSheet.absoluteFillObject,
		backgroundColor: "rgba(0,0,0,0.72)",
		alignItems: "center",
		justifyContent: "center",
		padding: 24,
	},
	card: {
		backgroundColor: "#1a1a1a",
		borderRadius: 16,
		borderWidth: 1,
		borderColor: "#2e2e2e",
		padding: 28,
		width: "100%",
		alignItems: "center",
	},
	title: {
		fontSize: 18,
		fontWeight: "600",
		color: "#f5f5f5",
		marginBottom: 6,
	},
	subtitle: {
		fontSize: 13,
		color: "#888",
		textAlign: "center",
		marginBottom: 24,
	},
	dotsRow: {
		flexDirection: "row",
		gap: 8,
	},
	dot: {
		width: 8,
		height: 8,
		borderRadius: 4,
		backgroundColor: "#4a9eff",
	},
	errorIcon: {
		fontSize: 36,
		color: "#e05a5a",
		marginBottom: 12,
	},
	successIcon: {
		fontSize: 36,
		color: "#4caf7d",
		marginBottom: 12,
	},
	errorBox: {
		backgroundColor: "#110e0e",
		borderWidth: 1,
		borderColor: "#3a1f1f",
		borderRadius: 8,
		paddingHorizontal: 14,
		paddingVertical: 10,
		width: "100%",
		marginBottom: 24,
	},
	errorMessage: {
		fontSize: 12,
		color: "#e05a5a",
		fontFamily: "monospace",
		lineHeight: 18,
	},
	retryButton: {
		backgroundColor: "#4a9eff",
		paddingHorizontal: 36,
		paddingVertical: 13,
		borderRadius: 10,
	},
	retryText: {
		color: "#fff",
		fontSize: 15,
		fontWeight: "600",
	},
});

const permStyles = StyleSheet.create({
	container: {
		flex: 1,
		backgroundColor: "#0f0f0f",
		alignItems: "center",
		justifyContent: "center",
		paddingHorizontal: 36,
	},
	icon: {
		fontSize: 48,
		marginBottom: 20,
	},
	title: {
		fontSize: 20,
		fontWeight: "600",
		color: "#f5f5f5",
		marginBottom: 10,
		textAlign: "center",
	},
	subtitle: {
		fontSize: 14,
		color: "#888",
		textAlign: "center",
		lineHeight: 22,
	},
});
