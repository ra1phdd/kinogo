import React, { useEffect } from 'react';
import noUiSlider, {API} from 'nouislider';

interface HTMLElementWithNoUiSlider extends HTMLElement {
    noUiSlider?: API;
}

export const useSlider = (sliderElement: HTMLElementWithNoUiSlider, minInput: HTMLInputElement, maxInput: HTMLInputElement, setFormData: React.Dispatch<React.SetStateAction<{
    genres: string[];
    year__min: string;
    year__max: string
}>>) => {
    useEffect(() => {
        const currentYear = new Date().getFullYear();

        if (sliderElement != null) {
            if (sliderElement.noUiSlider) {
                sliderElement.noUiSlider.destroy();
            }

            noUiSlider.create(sliderElement, {
                start: [1980, currentYear],
                connect: false,
                step: 1,
                range: {
                    min: 1980,
                    max: currentYear,
                },
            });

            sliderElement.noUiSlider?.on("update", (values: (string | number)[]) => {
                const min = parseInt(String(values[0]));
                const max = parseInt(String(values[1]));

                minInput.value = String(min);
                maxInput.value = String(max);

                setFormData(prevState => ({
                    ...prevState,
                    year__min: String(min),
                    year__max: String(max),
                }));
            });
        }
    }, [sliderElement, minInput, maxInput]);
};