import { createContext, useContext, ReactNode } from "react";

interface FormData {
    genres: string[];
    year__min: string;
    year__max: string;
}

const FormDataContext = createContext<FormData | null>(null);

export const useFormData = () => useContext(FormDataContext);

interface FormDataProviderProps {
    value: FormData;
    children: ReactNode;
}

export const FormDataProvider: React.FC<FormDataProviderProps> = ({ value, children }) => {
    return (
        <FormDataContext.Provider value={value}>
            {children}
        </FormDataContext.Provider>
    );
};
