import React, { useState, useEffect, useRef, KeyboardEvent } from 'react';

interface Option {
  value: string;
  label: string;
}

interface CustomSelectProps {
  options: Option[];
  onSelectChange: (selected: string[]) => void;
}

const CustomSelect: React.FC<CustomSelectProps> = ({ options, onSelectChange }) => {
  const [selectedOptions, setSelectedOptions] = useState<Option[]>([]);
  const [showOptions, setShowOptions] = useState(false);
  const [focusedOptionIndex, setFocusedOptionIndex] = useState<number>(-1);
  const selectBoxRef = useRef<HTMLDivElement>(null);
  const optionsRef = useRef<HTMLDivElement>(null);

  const handleSelect = (option: Option) => {
    let newSelectedOptions: Option[];
    if (selectedOptions.find(selectedOption => selectedOption.value === option.value)) {
      newSelectedOptions = selectedOptions.filter(selectedOption => selectedOption.value !== option.value);
    } else {
      newSelectedOptions = [...selectedOptions, option];
    }
    setSelectedOptions(newSelectedOptions);
    onSelectChange(newSelectedOptions.map(opt => opt.value));
    setFocusedOptionIndex(-1);
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLDivElement>) => {
    if (showOptions) {
      if (e.key === 'ArrowDown') {
        setFocusedOptionIndex((prevIndex) => (prevIndex + 1) % options.length);
      } else if (e.key === 'ArrowUp') {
        setFocusedOptionIndex((prevIndex) => (prevIndex - 1 + options.length) % options.length);
      } else if (e.key === 'Enter' && focusedOptionIndex >= 0) {
        handleSelect(options[focusedOptionIndex]);
      }
    }
  };

  const handleClick = () => {
    setShowOptions(prev => !prev);
    setFocusedOptionIndex(-1);
  };

  const handleOptionClick = (option: Option) => {
    if (showOptions) {
      handleSelect(option);
      setFocusedOptionIndex(-1); // Сбрасываем фокус при выборе опции
    }
  };

  const handleOutsideClick = (event: MouseEvent) => {
    if (selectBoxRef.current &&
        !selectBoxRef.current.contains(event.target as Node) &&
        optionsRef.current &&
        !optionsRef.current.contains(event.target as Node)) {
      setShowOptions(false);
    }
  };

  useEffect(() => {
    document.addEventListener('click', handleOutsideClick);

    return () => {
      document.removeEventListener('click', handleOutsideClick);
    };
  }, []);

  return (
      <div className="custom-select" tabIndex={0} onKeyDown={handleKeyDown}>
        <div
            className={`select-box ${showOptions ? 'clicked' : ''}`}
            ref={selectBoxRef}
            onClick={handleClick}
        >
          <p>{selectedOptions.length > 0
              ? selectedOptions.map(option => option.label).join(', ')
              : 'Выберите жанры'}</p>
          <span>▶</span>
        </div>
        {showOptions && (
            <div className={`options ${showOptions ? 'show' : ''}`} ref={optionsRef}>
              {options.map((option, index) => (
                  <div
                      key={option.value}
                      className={`option ${index === focusedOptionIndex ? 'focused' : ''} ${selectedOptions.find(selectedOption => selectedOption.value === option.value) ? 'selected' : ''}`}
                      onClick={() => handleOptionClick(option)}
                  >
                    {option.label}
                  </div>
              ))}
            </div>
        )}
      </div>
  );
};

export default CustomSelect;
