function СustomSelect(): void {
  const selectElement: HTMLSelectElement | null = document.querySelector('.custom-select');
  const selectWrapper: HTMLElement | null = document.querySelector('.select-wrapper');

  if (!selectElement || !selectWrapper) return;

  const options: HTMLOptionsCollection = selectElement.options;

  const selectResult: HTMLDivElement = document.createElement('div');
  selectResult.classList.add('select-result');
  selectWrapper.appendChild(selectResult);

  const firstDiv: HTMLDivElement = document.createElement('div');
  firstDiv.classList.add('select-view');
  selectResult.appendChild(firstDiv);

  const iconSelect: HTMLSpanElement = document.createElement('span');
  iconSelect.innerHTML = "▶";
  selectResult.appendChild(iconSelect);

  let firstOption: boolean = true;
  for (let i = 0; i < options.length; i++) {
    const option: HTMLOptionElement = options[i];
    if (firstOption) {
      firstDiv.innerHTML = option.innerHTML;
      firstOption = false;
    } else {
      const div: HTMLDivElement = document.createElement('div');
      const hr: HTMLHRElement = document.createElement('hr');
      div.innerHTML = option.innerHTML;
      div.dataset.value = option.value;
      div.classList.add('select-item');
      hr.classList.add('select-hr');
      selectResult.appendChild(div);
      selectResult.appendChild(hr);
    }
  }

  selectWrapper.addEventListener('click', (e: Event) => {
    selectResult.classList.toggle('select-active');
  });

  document.querySelectorAll('.select-result .select-item').forEach((item: Element) => {
    item.addEventListener('click', (e: Event) => {
      e.stopPropagation();
      if (selectResult.classList.contains('select-active')) {
        item.classList.toggle('selected');
        const selectedValues: string[] = [];
        document.querySelectorAll('.select-item.selected').forEach((selectedItem: Element) => {
          selectedValues.push((selectedItem as HTMLDivElement).dataset.value || '');
        });
        selectElement.value = selectedValues.join(',');
        for (let i = 0; i < options.length; i++) {
          const option: HTMLOptionElement = options[i];
          option.selected = selectedValues.indexOf(option.value) !== -1;
        }

        const selectedOptions: string[] = [];
        document.querySelectorAll('.select-item.selected').forEach((selectedItem: Element) => {
          selectedOptions.push((selectedItem as HTMLDivElement).innerHTML);
        });
        if (selectedOptions.length > 0) {
          firstDiv.innerHTML = selectedOptions.join(', ');
        } else {
          firstDiv.innerHTML = "Выберите жанр";
        }
      }
    });
  });

  document.addEventListener('click', (e: Event) => {
    if (!selectWrapper.contains(e.target as Node)) {
      selectResult.classList.remove('select-active');
    }
  });

  if (selectElement) {
    selectElement.addEventListener('change', (e: Event) => {
      document.querySelectorAll('.select-result div').forEach((item: Element) => {
        item.classList.remove('display-none');
      });
      const target: HTMLSelectElement = e.target as HTMLSelectElement;
      const selectedValue: string = target.options[target.selectedIndex].innerHTML;
      document.querySelectorAll('.select-result div').forEach((item: Element) => {
        if ((item as HTMLDivElement).innerHTML === selectedValue) {
          item.classList.add('display-none');
          firstDiv.style.display = 'block';
        }
      });
    });
  }
}

export default CustomSelect;
