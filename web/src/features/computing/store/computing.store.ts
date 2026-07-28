import { create } from "zustand";
import type { ComputeNode } from "../types/node.types";


interface ComputeStoreProps {
    error: string | null;
    currentNode: ComputeNode | null;
    nodes: ComputeNode[];

    setCurrentNode: (node: ComputeNode | null) => void;
    setNodes: (nodes: ComputeNode[]) => void;
    addNode: (node: ComputeNode) => void;
    updateNode: (node: ComputeNode) => void;
    deleteNode: (nodeId: string) => void;
    setError: (error: string | null) => void;
    reset: () => void;
}

const useComputeStore = create<ComputeStoreProps>((set) => ({
    currentNode: null,
    error: null,
    nodes: [],

    setCurrentNode: (node) => set({ currentNode: node }),

    setNodes: (nodes) => set({ nodes }),

    addNode: (node) =>
        set((state) => ({
            nodes: [...state.nodes, node],
            currentNode: node,
        })),

    updateNode: (node) =>
        set((state) => ({
            nodes: state.nodes.map((currentNode) =>
                currentNode.id === node.id ? node : currentNode,
            ),
            currentNode: state.currentNode?.id === node.id ? node : state.currentNode,
        })),

    deleteNode: (nodeId) =>
        set((state) => ({
            nodes: state.nodes.filter((node) => node.id !== nodeId),
            currentNode: state.currentNode?.id === nodeId ? null : state.currentNode,
        })),

    setError: (error) => set({ error }),

    reset: () =>
        set({
            currentNode: null,
            error: null,
            nodes: [],
        }),
}));

export default useComputeStore;
