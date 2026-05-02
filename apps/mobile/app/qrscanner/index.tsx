import React, { useState, useEffect } from "react";
import { CameraView, Camera, BarcodeScanningResult } from "expo-camera";
import { Text, View, StyleSheet, Button } from "react-native";

export default function QRScannerScreen() {
	const [hasPermission, setHasPermission] = useState<boolean | null>(null);
	const [scanned, setScanned] = useState(false);

	useEffect(() => {
		const getCameraPermissions = async () => {
			const { status } = await Camera.requestCameraPermissionsAsync();
			setHasPermission(status === "granted");
		};

		getCameraPermissions();
	}, []);

	const handleBarcodeScanned = ({ type, data }: BarcodeScanningResult) => {
		setScanned(true);
		alert(`Bar code with type ${type} and data ${data} has been scanned!`);
	};

	if (hasPermission === null) {
		return <Text>Requesting for camera permission</Text>;
	}
	if (hasPermission === false) {
		return <Text>No access to camera</Text>;
	}

	return (
		<View style={styles.container}>
			<CameraView
				onBarcodeScanned={scanned ? undefined : handleBarcodeScanned}
				barcodeScannerSettings={{
					barcodeTypes: ["qr", "pdf417"],
				}}
				style={StyleSheet.absoluteFillObject}
			/>
			{scanned && (
				<Button title={"Tap to Scan Again"} onPress={() => setScanned(false)} />
			)}
		</View>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		flexDirection: "column",
		justifyContent: "center",
	},
});

