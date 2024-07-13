export const formatDate = (timestamp: number): string => {
    const date = new Date(timestamp * 1000);

    const dateOptions: Intl.DateTimeFormatOptions = {
        year: 'numeric',
        month: 'numeric',
        day: 'numeric',
    };

    const timeOptions: Intl.DateTimeFormatOptions = {
        hour: 'numeric',
        minute: 'numeric',
    };

    return (
        date.toLocaleDateString(undefined, dateOptions) +
        ' ' +
        date.toLocaleTimeString(undefined, timeOptions)
    );
};
