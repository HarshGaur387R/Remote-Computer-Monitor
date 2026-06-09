import { useState, useEffect } from 'react';
import { View, Pressable, StyleSheet, TextInput, KeyboardAvoidingView, Platform, Alert, ScrollView } from 'react-native';
import { BlurView } from 'expo-blur';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { COLORS } from '@rcm/shared';
import { ComputerType, getTable, deleteRowById, insertRow, updateRowById, updateFieldsById } from '@/db/index';
import MaterialIcons from '@expo/vector-icons/MaterialIcons';

interface PcDetailedCardProps {
    computerId: number;
    onClose: () => void;
    onUpdateComputerName?: (id: number, newName: string) => Promise<void>;
    onDelete?: (id: number) => Promise<void>;
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

function validateFields(jsonResponse: any) {
    const systemInfo: SystemInfoType = {
        cpu_usage_percent: typeof jsonResponse.cpu_usage_percent === 'number' ? jsonResponse.cpu_usage_percent : 0,
        total_vram_mb: typeof jsonResponse.total_vram_mb === 'number' ? jsonResponse.total_vram_mb : 0,
        available_vram_mb: typeof jsonResponse.available_vram_mb === 'number' ? jsonResponse.available_vram_mb : 0,
        vram_used_percent: typeof jsonResponse.vram_used_percent === 'number' ? jsonResponse.vram_used_percent : 0,
        total_storage_gb: typeof jsonResponse.total_storage_gb === 'number' ? jsonResponse.total_storage_gb : 0,
        available_storage_gb: typeof jsonResponse.available_storage_gb === 'number' ? jsonResponse.available_storage_gb : 0,
        storage_used_percent: typeof jsonResponse.storage_used_percent === 'number' ? jsonResponse.storage_used_percent : 0,
    };

    return systemInfo
}

export function PcDetailedCard({ computerId, onClose, onUpdateComputerName, onDelete }: PcDetailedCardProps) {
    const [computer, setComputer] = useState<ComputerType | null>(null);
    const [isEditingName, setIsEditingName] = useState(false);
    const [editedName, setEditedName] = useState('');
    const [loading, setLoading] = useState(true);
    const [isSaving, setIsSaving] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    console.log("pc-card-details renderd")

    useEffect(() => {
        fetchComputerData();
    }, [computerId]);

    const fetchComputerData = async () => {
        const controller = new AbortController();
        const timeout = setTimeout(() => controller.abort(), 5000);
        let found: ComputerType | undefined;

        try {
            setLoading(true);
            const computers = await getTable();
            found = computers.find(c => c.id === computerId);

            if (found) {
                setEditedName(found.name || "Unknown device");
                setComputer(found); // Initialize with current state

                const response = await fetch(
                    `http://${found.LANIP}:${found.port}/system-info?auth=${found.authtoken}`,
                    {
                        signal: controller.signal,
                    }
                );

                if (!response.ok) {
                    throw new Error(`Agent responded with status ${response.status}`);
                }

                const jsonResponse = await response.json();
                const systemInfo = validateFields(jsonResponse);

                await updateFieldsById(found.id, { ...systemInfo, active: true });
                setComputer({ ...found, ...systemInfo, active: true });
            }
        } catch (error) {
            if (found) {
                const errorComputer = { ...found, active: false };
                setComputer(errorComputer);
                await updateFieldsById(found.id, { active: false });

                if (error instanceof Error && error.name === 'AbortError') {
                    console.error('Request timed out - setting device offline');
                } else {
                    console.error('Error fetching computer data:', error);
                }
            }
        } finally {
            clearTimeout(timeout);
            setLoading(false);
        }
    };

    const handleSaveName = async () => {
        if (!editedName.trim() || !computer) return;

        try {
            setIsSaving(true);
            if (onUpdateComputerName) {
                await onUpdateComputerName(computerId, editedName);
            }
            setComputer({ ...computer, name: editedName });
            setIsEditingName(false);
        } catch (error) {
            console.error('Error updating computer name:', error);
            setEditedName(computer?.name || "Unknown device");
        } finally {
            setIsSaving(false);
        }
    };

    const handleDelete = async () => {
        Alert.alert(
            'Delete Computer',
            `Are you sure you want to remove "${computer?.name || 'this computer'}"?`,
            [
                { text: 'Cancel', onPress: () => { }, style: 'cancel' },
                {
                    text: 'Delete',
                    onPress: async () => {
                        try {
                            setIsDeleting(true);
                            if (onDelete) {
                                await onDelete(computerId);
                            } else {
                                await deleteRowById(computerId);
                            }
                            onClose();
                        } catch (error) {
                            console.error('Error deleting computer:', error);
                            Alert.alert('Error', 'Failed to delete computer');
                        } finally {
                            setIsDeleting(false);
                        }
                    },
                    style: 'destructive',
                },
            ]
        );
    };

    if (!computer && !loading) {
        return null;
    }

    return (
        <BlurView style={styles.blurContainer} intensity={10}>
            <KeyboardAvoidingView
                behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
                style={styles.overlay}
            >
                <Pressable style={styles.dismissArea} onPress={onClose} />

                <ThemedView style={styles.cardContainer}>
                    {/* Close Button */}
                    <Pressable
                        style={styles.closeButton}
                        onPress={onClose}
                        hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
                    >
                        <MaterialIcons name="close" size={24} color="#888" />
                    </Pressable>

                    {/* Header Section */}
                    <View style={styles.headerSection}>
                        {isEditingName ? (
                            <View style={styles.editNameContainer}>
                                <TextInput
                                    style={styles.nameInput}
                                    value={editedName}
                                    onChangeText={setEditedName}
                                    placeholder="Enter computer name"
                                    placeholderTextColor="#666"
                                    autoFocus
                                    editable={!isSaving}
                                />
                                <Pressable
                                    style={[styles.saveButton, isSaving && styles.savingButton]}
                                    onPress={handleSaveName}
                                    disabled={isSaving}
                                >
                                    <MaterialIcons
                                        name={isSaving ? 'hourglass-bottom' : 'check'}
                                        size={16}
                                        color={isSaving ? '#666' : '#fff'}
                                    />
                                </Pressable>
                                <Pressable
                                    style={styles.cancelButton}
                                    onPress={() => {
                                        setIsEditingName(false);
                                        setEditedName(computer?.name || "Unknown device");
                                    }}
                                    disabled={isSaving}
                                >
                                    <MaterialIcons name="close" size={16} color="#888" />
                                </Pressable>
                            </View>
                        ) : (
                            <View style={styles.nameRow}>
                                <ThemedText style={styles.computerName}>
                                    {computer?.name || 'Unnamed Device'}
                                </ThemedText>
                                <Pressable
                                    style={styles.editButton}
                                    onPress={() => setIsEditingName(true)}
                                    hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                                >
                                    {/* <MaterialIcons name="edit" size={18} color={COLORS.primary} /> */}
                                </Pressable>
                            </View>
                        )}

                        <View style={styles.statusRow}>
                            <View
                                style={[
                                    styles.statusDot,
                                    { backgroundColor: computer?.active ? '#4ade80' : '#ef4444' }
                                ]}
                            />
                            <ThemedText style={styles.statusText}>
                                {computer?.active ? 'Online' : 'Offline'}
                            </ThemedText>
                        </View>
                    </View>

                    <View style={styles.divider} />

                    {/* Details Section */}
                    <ScrollView>

                        <View style={styles.detailsSection}>
                            {loading ? (
                                <ThemedText style={styles.loadingText}>Loading details...</ThemedText>
                            ) : (
                                <>
                                    <DetailRow label="IP Address" value={computer?.LANIP || 'N/A'} />
                                    <DetailRow label="Port" value={computer?.port?.toString() || 'N/A'} />
                                    <DetailRow
                                        label="Auth Token"
                                        value={computer?.authtoken ? `${computer.authtoken.substring(0, 10)}...` : 'N/A'}
                                        isCopyable
                                        fullValue={computer?.authtoken}
                                    />
                                    <DetailRow label="Device ID" value={computer?.id?.toString() || 'N/A'} />

                                    {/* System Stats */}
                                    <DetailRow
                                        label="CPU Usage"
                                        value={computer?.cpu_usage_percent ? `${computer.cpu_usage_percent.toFixed(1)}%` : 'N/A'}
                                    />
                                    <DetailRow
                                        label="RAM Total"
                                        value={computer?.total_vram_mb ? `${computer.total_vram_mb} MB` : 'N/A'}
                                    />
                                    <DetailRow
                                        label="RAM Available"
                                        value={computer?.available_vram_mb ? `${computer.available_vram_mb} MB` : 'N/A'}
                                    />
                                    <DetailRow
                                        label="RAM Usage"
                                        value={computer?.vram_used_percent ? `${computer.vram_used_percent.toFixed(1)}%` : 'N/A'}
                                    />
                                    <DetailRow
                                        label="Storage Total"
                                        value={computer?.total_storage_gb ? `${computer.total_storage_gb} GB` : 'N/A'}
                                    />
                                    <DetailRow
                                        label="Storage Available"
                                        value={computer?.available_storage_gb ? `${computer.available_storage_gb} GB` : 'N/A'}
                                    />
                                    <DetailRow
                                        label="Storage Usage"
                                        value={computer?.storage_used_percent ? `${computer.storage_used_percent.toFixed(1)}%` : 'N/A'}
                                    />
                                </>
                            )}
                        </View>
                    </ScrollView>

                    {/* Action Buttons */}
                    <View style={styles.buttonRow}>
                        <Pressable style={[styles.actionButton, styles.connectButton]}>
                            <ThemedText style={styles.buttonText}>Connect</ThemedText>
                        </Pressable>
                        <Pressable
                            style={[styles.actionButton, styles.deleteButton]}
                            onPress={handleDelete}
                            disabled={isDeleting}
                        >
                            <ThemedText style={styles.deleteButtonText}>
                                {isDeleting ? 'Removing...' : 'Remove'}
                            </ThemedText>
                        </Pressable>
                    </View>
                </ThemedView>
            </KeyboardAvoidingView>
        </BlurView>
    );
}

interface DetailRowProps {
    label: string;
    value: string;
    isCopyable?: boolean;
    fullValue?: string;
}

function DetailRow({ label, value, isCopyable = false, fullValue }: DetailRowProps) {
    const [copied, setCopied] = useState(false);

    const handleCopy = () => {
        if (fullValue) {
            // In a real app, use Clipboard.setString(fullValue)
            setCopied(true);
            setTimeout(() => setCopied(false), 1500);
        }
    };

    return (
        <View style={styles.detailRow}>
            <ThemedText style={styles.detailLabel}>{label}</ThemedText>
            <Pressable
                style={styles.detailValueContainer}
                onPress={isCopyable ? handleCopy : undefined}
            >
                <ThemedText style={styles.detailValue} numberOfLines={1}>
                    {value}
                </ThemedText>
                {isCopyable && (
                    <MaterialIcons
                        name={copied ? 'done' : 'content-copy'}
                        size={14}
                        color={copied ? '#4ade80' : '#888'}
                        style={styles.copyIcon}
                    />
                )}
            </Pressable>
        </View>
    );
}

const styles = StyleSheet.create({
    blurContainer: {
        ...StyleSheet.absoluteFillObject,
        justifyContent: 'center',
        alignItems: 'center',
        zIndex: 1000,
    },
    overlay: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        width: '100%',
        height: '100%',
    },
    dismissArea: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
    },
    cardContainer: {
        width: '85%',
        maxHeight: '80%',
        backgroundColor: '#0C0C0C',
        borderRadius: 16,
        borderWidth: 1,
        borderColor: 'rgba(15, 40, 84, 0.5)',  // fixed: was invalid 'rgb(15, 40, 84, 0.5)'
        padding: 20,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 8 },
        shadowOpacity: 0.5,
        shadowRadius: 12,
        elevation: 12,
        zIndex: 1001,
    },
    closeButton: {
        position: 'absolute',
        top: 12,
        right: 12,
        padding: 8,
        zIndex: 10,
    },
    headerSection: {
        marginTop: 4,
    },
    nameRow: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        marginBottom: 8,
    },
    computerName: {
        fontSize: 18,
        fontWeight: '700',
        color: 'white',
        flex: 1,
    },
    editButton: {
        padding: 4,
        marginLeft: 8,
    },
    editNameContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: 8,
        gap: 8,
    },
    nameInput: {
        flex: 1,
        backgroundColor: 'rgba(15, 40, 84, 0.3)',  // fixed: was invalid 'rgb(15, 40, 84, 0.3)'
        borderWidth: 1,
        borderColor: COLORS.primary,
        borderRadius: 8,
        paddingHorizontal: 12,
        paddingVertical: 8,
        color: 'white',
        fontSize: 16,
        fontWeight: '600',
    },
    saveButton: {
        backgroundColor: COLORS.primary,
        borderRadius: 8,
        padding: 8,
        justifyContent: 'center',
        alignItems: 'center',
    },
    savingButton: {
        backgroundColor: 'rgba(15, 40, 84, 0.5)',  // fixed: was invalid 'rgba(15, 40, 84, 0.5)'
    },
    cancelButton: {
        backgroundColor: 'rgba(15, 40, 84, 0.3)',  // fixed: was invalid 'rgb(15, 40, 84, 0.3)'
        borderRadius: 8,
        padding: 8,
        justifyContent: 'center',
        alignItems: 'center',
        borderWidth: 1,
        borderColor: 'rgba(15, 40, 84, 0.5)',  // fixed: was invalid 'rgb(15, 40, 84, 0.5)'
    },
    statusRow: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 6,
    },
    statusDot: {
        width: 8,
        height: 8,
        borderRadius: 4,
    },
    statusText: {
        fontSize: 12,
        color: '#888',
        fontWeight: '500',
    },
    divider: {
        height: 1,
        backgroundColor: 'rgba(15, 40, 84, 0.5)',  // fixed: was invalid 'rgb(15, 40, 84, 0.5)'
        marginVertical: 16,
    },
    detailsSection: {
        gap: 12,
    },
    loadingText: {
        textAlign: 'center',
        color: '#888',
        fontSize: 14,
        paddingVertical: 12,
    },
    detailRow: {
        gap: 8,
    },
    detailLabel: {
        fontSize: 11,
        color: '#888',
        fontWeight: '500',
        marginBottom: 2,
    },
    detailValueContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        backgroundColor: 'rgba(15, 40, 84, 0.2)',  // fixed: was invalid 'rgb(15, 40, 84, 0.2)'
        borderRadius: 6,
        paddingHorizontal: 10,
        paddingVertical: 8,
    },
    detailValue: {
        fontSize: 12,
        color: '#ccc',
        fontWeight: '500',
        flex: 1,
    },
    copyIcon: {
        marginLeft: 8,
    },
    buttonRow: {
        flexDirection: 'row',
        gap: 10,
        marginTop: 16,
    },
    actionButton: {
        flex: 1,
        paddingVertical: 10,
        borderRadius: 8,
        justifyContent: 'center',
        alignItems: 'center',
    },
    connectButton: {
        backgroundColor: COLORS.primary,
    },
    deleteButton: {
        backgroundColor: 'rgba(239, 68, 68, 0.2)',
        borderWidth: 1,
        borderColor: 'rgba(239, 68, 68, 0.5)',
    },
    buttonText: {
        fontSize: 14,
        fontWeight: '600',
        color: 'white',
    },
    deleteButtonText: {
        fontSize: 14,
        fontWeight: '600',
        color: '#ef4444',
    },
});