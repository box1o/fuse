export const normalizeDeviceCode = (value: string): string => {
    return value.toUpperCase().replace(/[^A-Z0-9]/g, "").slice(0, 8);
};

export const formatDeviceCode = (value: string): string => {
    const normalized = normalizeDeviceCode(value);
    if (normalized.length <= 4) return normalized;
    return `${normalized.slice(0, 4)}-${normalized.slice(4)}`;
};
