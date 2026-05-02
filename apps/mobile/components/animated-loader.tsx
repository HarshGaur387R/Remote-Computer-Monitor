// components/AnimatedLoader.tsx
import React from 'react';
import { View, StyleSheet } from 'react-native';
import { WebView } from 'react-native-webview';
import { loader1Xml } from './loaders/loader1';
import { loader2Xml } from './loaders/loader2';

const LOADERS = {
  loader1: loader1Xml,
  loader2: loader2Xml,
} as const;

export type LoaderType = keyof typeof LOADERS;

interface AnimatedLoaderProps {
  type?: LoaderType;
  size?: number;
  backgroundColor?: string;
}

export default function AnimatedLoader({
  type = 'loader1',
  size = 100,
  backgroundColor = 'transparent',
}: AnimatedLoaderProps) {
  const svgXml = LOADERS[type];

  const html = `
    <!DOCTYPE html>
    <html>
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <style>
          * { margin: 0; padding: 0; box-sizing: border-box; }
          html, body {
            width: 100%;
            height: 100%;
            background: transparent;
            display: flex;
            align-items: center;
            justify-content: center;
            overflow: hidden;
          }
          svg {
            width: 100%;
            height: 100%;
          }
        </style>
      </head>
      <body>
        ${svgXml}
      </body>
    </html>
  `;

  return (
    <View style={[styles.container, { width: size, height: size }]}>
      <WebView
        source={{ html }}
        style={styles.webview}
        scrollEnabled={false}
        showsHorizontalScrollIndicator={false}
        showsVerticalScrollIndicator={false}
        backgroundColor={backgroundColor}
        originWhitelist={['*']}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    overflow: 'hidden',
  },
  webview: {
    flex: 1,
    backgroundColor: 'transparent',
  },
});
