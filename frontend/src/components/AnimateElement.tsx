function AnimateElement(element: HTMLElement, animation: string, delay: number) {
    setTimeout(function () {
        if (element == document.querySelector(".section__movie")) {
            element.style.display = "flex";
        } else element.style.display = "block";
        element.classList.add("animate__animated", animation, "animate__faster");
    }, delay);
}

export default AnimateElement;