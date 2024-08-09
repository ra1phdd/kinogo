import React, {useState} from "react";

export const useTextareaState = (initialValue: string = '') => {
    const [value, setValue] = useState<string>(initialValue);

    const handleChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
        setValue(event.target.value);
    };

    const reset = () => {
        setValue(initialValue);
    };

    return {
        value,
        onChange: handleChange,
        reset,
    };
}