import {useEffect} from "react";
import AnimateElement from "@components/AnimateElement.tsx";

function SearchAside() {
    useEffect(() => {
        const asideSearch = document.querySelector<HTMLElement>(".aside__search");

        if (asideSearch !== null) {
            AnimateElement(asideSearch, "animate__fadeInRight", 150);
        }
    }, []);

    return (
        <div className="aside__search">
            <h3>Поиск</h3>
            <form action="/search" method="POST">
                <input type="text" name="search" id="aside__search-input" placeholder="Введите запрос"/>
                <button type="submit"><img src="/src/assets/images/search.svg" alt="Поиск"/></button>
            </form>
            <div className="aside__search-results" id="aside__search-results"></div>
        </div>
    )
}

export default SearchAside
