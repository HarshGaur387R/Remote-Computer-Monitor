
import { ViewStyle } from 'react-native';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { IconSymbol } from '@/components/ui/icon-symbol';

type NetworkErrorScreenTypes = {
	error: string;
	message: string;
	container_style: ViewStyle;
}

function NetworkErrorScreen({ error, message, container_style }: NetworkErrorScreenTypes) {
	return (
		<ThemedView style={container_style}>
			<IconSymbol name='network.slash' size={130} color={"#1C4D8D"} />
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
	NetworkErrorScreen,
	NetworkErrorScreenTypes

}
