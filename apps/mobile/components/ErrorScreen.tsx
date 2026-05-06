import { ViewStyle } from 'react-native';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { IconSymbol } from '@/components/ui/icon-symbol';

type ErrorScreenTypes = {
	error: string;
	message: string;
	container_style: ViewStyle;
}

function ErrorScreen({ error, message, container_style }: ErrorScreenTypes) {
	return (
		<ThemedView style={container_style}>
			<IconSymbol name='stop' size={130} color={"#1C4D8D"} />
			<ThemedText style={{
				color: "#1C4D8D",
				fontWeight: '700',
				fontSize: 32,
				textAlign: 'center',
				padding: 4, 
				paddingBottom: 20
			}}>
				{error}
			</ThemedText>
			<ThemedText style={{
				color: 'darkgrey',
				fontWeight: '600',
				fontSize: 20,
				textAlign: 'center',
				lineHeight: 24
			}
			}>{message}</ThemedText>
		</ThemedView>
	)
}

export {
	ErrorScreen,
	ErrorScreenTypes

}
