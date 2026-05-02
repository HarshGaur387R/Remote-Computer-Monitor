import { Tabs } from 'expo-router';
import React from 'react';
import { HapticTab } from '@/components/haptic-tab';
import { IconSymbol } from '@/components/ui/icon-symbol';

export default function TabLayout() {

  return (
    <Tabs
      screenOptions={{
        tabBarInactiveTintColor: "#2c5b87",
        tabBarActiveTintColor: "#4988C4",
        headerShown: false,
        tabBarButton: HapticTab,
      }}>
      <Tabs.Screen
        name="index"
        options={{
          title: 'Home',
          tabBarIcon: ({ color }) => <IconSymbol size={28} name="house.fill" color={color} />,
        }}
      />
      <Tabs.Screen
        name="computers"
        options={{
          title: 'Computers',
          tabBarIcon: ({ color }) => <IconSymbol size={28} name="pc" color={color} />,
        }}
      />
    </Tabs>
  );
}
