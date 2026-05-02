import { ViewStyle } from 'react-native';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

type LoadingScreenType = {
	container_style: ViewStyle;
	animated_loader_child: React.ReactNode;
	text: string;
}

function Loading_screen({ container_style, animated_loader_child, text }: LoadingScreenType) {
	return (
		<ThemedView style={container_style}>
			{animated_loader_child}
			<ThemedText style={{ color: 'darkgrey', fontSize: 18 }} >{text}</ThemedText>
		</ThemedView>
	)
}

export {Loading_screen, LoadingScreenType}
