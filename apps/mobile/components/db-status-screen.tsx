import React, { useEffect, useRef } from "react";
import {
	View,
	Text,
	StyleSheet,
	Animated,
	TouchableOpacity,
	Easing,
} from "react-native";

// ─── Loading Screen ───────────────────────────────────────────────

export function DBLoadingScreen() {
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

		a1.start();
		a2.start();
		a3.start();

		return () => {
			a1.stop();
			a2.stop();
			a3.stop();
		};
	}, []);

	const dotStyle = (anim: Animated.Value) => ({
		opacity: anim,
		transform: [
			{
				translateY: anim.interpolate({
					inputRange: [0, 1],
					outputRange: [0, -6],
				}),
			},
		],
	});

	return (
		<View style={styles.container}>
			<Text style={styles.title}>Initializing</Text>
			<Text style={styles.subtitle}>Setting up local database…</Text>
			<View style={styles.dotsRow}>
				<Animated.View style={[styles.dot, dotStyle(dot1)]} />
				<Animated.View style={[styles.dot, dotStyle(dot2)]} />
				<Animated.View style={[styles.dot, dotStyle(dot3)]} />
			</View>
		</View>
	);
}

// ─── Error Screen ─────────────────────────────────────────────────

interface DBErrorScreenProps {
	error: Error;
	onRetry: () => void;
}

export function DBErrorScreen({ error, onRetry }: DBErrorScreenProps) {
	return (
		<View style={styles.container}>
			<View style={styles.errorIconContainer}>
				<Text style={styles.errorIcon}>⚠</Text>
			</View>
			<Text style={styles.title}>Database Error</Text>
			<Text style={styles.subtitle}>Could not initialize the local database.</Text>
			<View style={styles.errorBox}>
				<Text style={styles.errorMessage} numberOfLines={4} ellipsizeMode="tail">
					{error.message}
				</Text>
			</View>
			<TouchableOpacity style={styles.retryButton} onPress={onRetry} activeOpacity={0.75}>
				<Text style={styles.retryText}>Try Again</Text>
			</TouchableOpacity>
		</View>
	);
}

// ─── Styles ───────────────────────────────────────────────────────

const styles = StyleSheet.create({
	container: {
		flex: 1,
		backgroundColor: "#0f0f0f",
		alignItems: "center",
		justifyContent: "center",
		paddingHorizontal: 32,
	},

	// shared
	title: {
		fontSize: 22,
		fontWeight: "600",
		color: "#f5f5f5",
		marginBottom: 8,
	},
	subtitle: {
		fontSize: 14,
		color: "#888",
		textAlign: "center",
		marginBottom: 32,
	},

	// loading dots
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

	// error
	errorIconContainer: {
		marginBottom: 16,
	},
	errorIcon: {
		fontSize: 40,
		color: "#e05a5a",
	},
	errorBox: {
		backgroundColor: "#1a1a1a",
		borderWidth: 1,
		borderColor: "#2e2e2e",
		borderRadius: 8,
		paddingHorizontal: 16,
		paddingVertical: 12,
		width: "100%",
		marginBottom: 28,
	},
	errorMessage: {
		fontSize: 13,
		color: "#e05a5a",
		fontFamily: "monospace",
		lineHeight: 20,
	},
	retryButton: {
		backgroundColor: "#4a9eff",
		paddingHorizontal: 40,
		paddingVertical: 14,
		borderRadius: 10,
	},
	retryText: {
		color: "#fff",
		fontSize: 15,
		fontWeight: "600",
	},
});
