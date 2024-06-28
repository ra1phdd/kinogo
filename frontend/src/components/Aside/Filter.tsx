import '../../assets/css/src/nouislider.css';
import noUiSlider from 'nouislider';
import customSelect from '../../assets/js/src/custom-select';
import React, {useEffect, useState} from "react";
import AnimateElement from "@components/AnimateElement.tsx";
import { useNavigate } from 'react-router-dom';

const FilterAside = React.memo(() => {
    const [formData, setFormData] = useState<{ genre: string[], year__min: string, year__max: string }>({
        genre: [],
        year__min: '',
        year__max: ''
    });
    const navigate = useNavigate();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const { name, value } = e.target;

        if (name === 'genre') {
            const selectedOptions = (e.target as HTMLSelectElement).selectedOptions;
            const selectedGenres = Array.from(selectedOptions).map(option => option.value);
            setFormData(prevState => ({
                ...prevState,
                genre: selectedGenres,
            }));
        } else {
            setFormData({
                ...formData,
                [name]: value,
            });
        }
    };

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        console.log(formData)
        e.preventDefault();
        navigate('/filter', { state: formData });
    };

    useEffect(() => {
        const asideFilter = document.querySelector<HTMLElement>(".aside__filter");

        if (asideFilter !== null) {
            AnimateElement(asideFilter, "animate__fadeInRight", 300);
        }

        customSelect();

        const slider: any = document.getElementById("slider");
        const sliderValueMin: any = document.getElementById("slider-min");
        const sliderValueMax: any = document.getElementById("slider-max");

        if (slider != null) {
            if (slider.noUiSlider) {
                slider.noUiSlider.destroy();
            }

            noUiSlider.create(slider, {
                start: [1980, 2024],
                connect: false,
                step: 1,
                range: {
                    min: 1980,
                    max: 2024,
                },
            });

            slider.noUiSlider.on("update", function (values: any) {
                sliderValueMin.value = parseInt(values[0]);
                sliderValueMax.value = parseInt(values[1]);
            });

            sliderValueMin.addEventListener("input", function (this: HTMLInputElement) {
                const value = parseInt(this.value);
                if (!isNaN(value) && value > 1980) {
                    slider.noUiSlider.set([value, null]);
                }
            });

            sliderValueMax.addEventListener("input", function (this: HTMLInputElement) {
                const value = parseInt(this.value);
                if (!isNaN(value) && value < 2024) {
                    slider.noUiSlider.set([null, value]);
                }
            });

            sliderValueMin.addEventListener("change", function (this: HTMLInputElement) {
                slider.noUiSlider.set([this.value, null]);
            });

            sliderValueMax.addEventListener("change", function (this: HTMLInputElement) {
                slider.noUiSlider.set([null, this.value]);
            });
        }
    }, []);

    return (
        <div className="aside__filter">
            <h3>Сортировка</h3>
            <form onSubmit={handleSubmit}>
                <select className="custom-select" name="genre" multiple onChange={handleChange}>
                    <option value="выбрать">Выберите жанр</option>
                    <option value="аниме">Аниме</option>
                    <option value="биография">Биография</option>
                    <option value="боевик">Боевик</option>
                    <option value="вестерн">Вестерн</option>
                    <option value="военный">Военный</option>
                    <option value="детектив">Детектив</option>
                    <option value="детский">Детский</option>
                    <option value="документальный">Документальный</option>
                    <option value="драма">Драма</option>
                    <option value="игра">Игра</option>
                    <option value="история">История</option>
                    <option value="комедия">Комедия</option>
                    <option value="концерт">Концерт</option>
                    <option value="короткометражка">Короткометражка</option>
                    <option value="криминал">Криминал</option>
                    <option value="мелодрама">Мелодрама</option>
                    <option value="музыка">Музыка</option>
                    <option value="мультфильм">Мультфильм</option>
                    <option value="мюзикл">Мюзикл</option>
                    <option value="новости">Новости</option>
                    <option value="приключения">Приключения</option>
                    <option value="семейный">Семейный</option>
                    <option value="спорт">Спорт</option>
                    <option value="ток-шоу">Ток-шоу</option>
                    <option value="триллер">Триллер</option>
                    <option value="ужасы">Ужасы</option>
                    <option value="фантастика">Фантастика</option>
                    <option value="фэнтези">Фэнтези</option>
                </select>
                <div className="select-wrapper"></div>
                <div className="slider__year">
                    <h4>Выберите год</h4>
                    <div id="slider"></div>
                    <div className="slider__value">
                        <input type="text" id="slider-min" name="year__min" value={formData.year__min}
                               onChange={handleChange}/>
                        <input type="text" id="slider-max" name="year__max" value={formData.year__max}
                               onChange={handleChange}/>
                    </div>
                </div>
                <button type="submit">
                    <div id="circle"></div>
                    <a href="#">Отправить</a>
                </button>
            </form>
        </div>
    )
});

export default FilterAside
