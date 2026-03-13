import { StyleSheet } from 'react-native';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { SafeAreaView } from 'react-native-safe-area-context';
import { APP_NAME} from "@rcm/shared"

export default function HomeScreen() {
	return (
		<SafeAreaView style={styles.container}>
			<ThemedView style={styles.mainScreen}>
				<ThemedText style={styles.title}>
				{APP_NAME}
				</ThemedText>
			</ThemedView>
		</SafeAreaView>
	);
}

const styles = StyleSheet.create({
	container: {
		padding: 5,
		width: "100%",
		height: "100%",
		flex: 1,
		alignItems: 'center',
		justifyContent: 'center',
		backgroundColor: "#0C0C0C"
	},
	mainScreen: {
		padding: 10,
		backgroundColor: '#481E14',
		borderRadius: 16,
		width: "90%",
		height: 600,
		borderWidth : 2,
		borderColor : "#F2613F"
	},
	title: {
		fontSize: 28,
		color: 'white',
		textAlign: 'center',
		padding: 5,
		fontWeight: 'semibold',
		lineHeight: 30,
	}
});
