import { DBLoadingScreen, DBErrorScreen } from "@/components/db-status-screen"
import { DarkTheme, DefaultTheme, ThemeProvider } from '@react-navigation/native';
import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import 'react-native-reanimated';

import { useColorScheme } from '@/hooks/use-color-scheme';
import { initDB } from '@/db';
import { useCallback, useEffect, useState } from 'react';

export const unstable_settings = {
	anchor: '(tabs)',
};

export default function RootLayout() {
	const colorScheme = useColorScheme();
	const [dbReady, setDbReady] = useState(false);
	const [dbError, setDbError] = useState<Error | null>(null);

	const initializeDB = useCallback(() => {
		setDbError(null);
		setDbReady(false);
		initDB()
			.then(() => setDbReady(true))
			.catch((err) => setDbError(err instanceof Error ? err : new Error(String(err))));
	}, []);

	useEffect(() => {
		initializeDB();
	}, []);

	if (dbError) return <DBErrorScreen error={dbError} onRetry={initializeDB} />;
	if (!dbReady) return <DBLoadingScreen />;

	return (
		<ThemeProvider value={colorScheme === "dark" ? DarkTheme : DefaultTheme}>
			<Stack>
				<Stack.Screen name="(tabs)" options={{ headerShown: false }} />
				<Stack.Screen name="qrscanner/index" options={{ headerShown: true, title: "Scan QR" }} />
				<Stack.Screen name="modal" options={{ presentation: "modal", title: "Modal" }} />
			</Stack>
			<StatusBar style="auto" />
		</ThemeProvider>
	);
}
