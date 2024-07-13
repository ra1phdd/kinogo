import '@assets/styles/vendor/nouislider.min.css';
import noUiSlider, { API } from 'nouislider';
import CustomSelect from '@components/specific/CustomSelect.tsx';
import React, { useEffect, useState } from "react";
import AnimateElement from "@components/common/AnimateElement.tsx";
import { useNavigate } from "react-router-dom";

interface HTMLElementWithNoUiSlider extends HTMLElement {
    noUiSlider?: API;
}
const FilterAside = React.memo(() => {
    const [formData, setFormData] = useState<{ genres: string[], year__min: string, year__max: string }>({
        genres: [],
        year__min: '',
        year__max: ''
    });
    const navigate = useNavigate();

    const handleSelectChange = (selected: string[]) => {
        setFormData(prevState => ({
            ...prevState,
            genres: selected
        }));
    };

    const handleMinChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = parseInt(e.target.value);
        if (!isNaN(value) && value >= 1980) {
            setFormData(prevState => ({
                ...prevState,
                year__min: value.toString(),
            }));
            const slider = document.getElementById("slider") as HTMLElementWithNoUiSlider;
            if (slider && slider.noUiSlider) {
                slider.noUiSlider.set([value, null]);
            }
        }
    };

    const handleMaxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = parseInt(e.target.value);
        if (!isNaN(value) && value <= new Date().getFullYear()) {
            setFormData(prevState => ({
                ...prevState,
                year__max: value.toString(),
            }));
            const slider = document.getElementById("slider") as HTMLElementWithNoUiSlider;
            if (slider && slider.noUiSlider) {
                slider.noUiSlider.set([null, value]);
            }
        }
    };

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        navigate('/filter', { state: formData });

        window.location.reload();
    };

    useEffect(() => {
        const asideFilter = document.querySelector<HTMLElement>(".aside__filter");
        const currentYear = new Date().getFullYear();

        if (asideFilter !== null) {
            AnimateElement(asideFilter, "animate__fadeInRight", 75);
        }

        const slider = document.getElementById("slider") as HTMLElementWithNoUiSlider;
        const sliderValueMin = document.getElementById("slider-min") as HTMLInputElement;
        const sliderValueMax = document.getElementById("slider-max") as HTMLInputElement;

        if (slider != null) {
            if (slider.noUiSlider) {
                slider.noUiSlider.destroy();
            }

            noUiSlider.create(slider, {
                start: [1980, currentYear],
                connect: false,
                step: 1,
                range: {
                    min: 1980,
                    max: currentYear,
                },
            });

            slider.noUiSlider?.on("update", (values: (string | number)[]) => {
                const min = parseInt(String(values[0]));
                const max = parseInt(String(values[1]));

                sliderValueMin.value = String(min);
                sliderValueMax.value = String(max);

                setFormData(prevState => ({
                    ...prevState,
                    year__min: String(min),
                    year__max: String(max),
                }));
            });
        }
    }, []);

    const genreOptions = [
        { value: 'аниме', label: 'Аниме' },
        { value: 'биография', label: 'Биография' },
        { value: 'боевик', label: 'Боевик' },
        { value: 'вестерн', label: 'Вестерн' },
        { value: 'военный', label: 'Военный' },
        { value: 'детектив', label: 'Детектив' },
        { value: 'детский', label: 'Детский' },
        { value: 'документальный', label: 'Документальный' },
        { value: 'драма', label: 'Драма' },
        { value: 'игра', label: 'Игра' },
        { value: 'история', label: 'История' },
        { value: 'комедия', label: 'Комедия' },
        { value: 'концерт', label: 'Концерт' },
        { value: 'короткометражка', label: 'Короткометражка' },
        { value: 'криминал', label: 'Криминал' },
        { value: 'мелодрама', label: 'Мелодрама' },
        { value: 'музыка', label: 'Музыка' },
        { value: 'мультфильм', label: 'Мультфильм' },
        { value: 'мюзикл', label: 'Мюзикл' },
        { value: 'новости', label: 'Новости' },
        { value: 'приключения', label: 'Приключения' },
        { value: 'семейный', label: 'Семейный' },
        { value: 'спорт', label: 'Спорт' },
        { value: 'ток-шоу', label: 'Ток-шоу' },
        { value: 'триллер', label: 'Триллер' },
        { value: 'ужасы', label: 'Ужасы' },
        { value: 'фантастика', label: 'Фантастика' },
        { value: 'фэнтези', label: 'Фэнтези' },
    ];

    return (
        <div className="aside__filter">
            <h3>Сортировка</h3>
            <form onSubmit={handleSubmit}>
                <CustomSelect options={genreOptions} onSelectChange={handleSelectChange} />
                <div className="slider__year">
                    <h4>Выберите год</h4>
                    <div id="slider"></div>
                    <div className="slider__value">
                        <input type="text" id="slider-min" name="year__min" value={formData.year__min}
                               onChange={handleMinChange}/>
                        <input type="text" id="slider-max" name="year__max" value={formData.year__max}
                               onChange={handleMaxChange}/>
                    </div>
                </div>
                <button type="submit">
                    <div id="circle"></div>
                    <span>Отправить</span>
                </button>
            </form>
        </div>
    )
});

export default FilterAside;