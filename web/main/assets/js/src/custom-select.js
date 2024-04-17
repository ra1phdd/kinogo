"use strict";

var option = document.querySelectorAll('.custom-select option'),
selectWrapper = document.querySelector('.select-wrapper');

var selectResult = document.createElement('div');
selectResult.classList.add('select-result');
selectWrapper.appendChild(selectResult);

var firstDiv = document.createElement('div');
firstDiv.classList.add('select-view');
selectResult.appendChild(firstDiv);

var iconSelect = document.createElement('span');
iconSelect.innerHTML = "▶";
selectResult.appendChild(iconSelect);

var firstOption = true;
option.forEach(function (item) {
  if (firstOption) {
    firstDiv.innerHTML = item.innerHTML;
    firstOption = false;
  } else{
    var div = document.createElement('div');
    var hr = document.createElement('hr');
    div.innerHTML = item.innerHTML;
    div.dataset.value = item.value;
    div.classList.add('select-item');
    hr.classList.add('select-hr');
    selectResult.appendChild(div);
    selectResult.appendChild(hr);
  }
});

document.querySelector('.select-wrapper').addEventListener('click', function (e) {
  selectResult.classList.toggle('select-active');
});

document.querySelectorAll('.select-result .select-item').forEach(function (item) {
  item.addEventListener('click', function (e) {
    e.stopPropagation();
    if (selectResult.classList.contains('select-active')) {
      item.classList.toggle('selected');
      var selectedValues = Array.from(document.querySelectorAll('.select-item.selected')).map(function(item) {
        return item.dataset.value;
      });
      document.querySelector('.custom-select').value = selectedValues;
      option.forEach(function (item) {
        item.selected = selectedValues.includes(item.value);
      });

      var selectedOptions = Array.from(document.querySelectorAll('.select-item.selected')).map(function(item) {
        return item.innerHTML;
      });
      if (selectedOptions.length > 0) {
        document.querySelector('.select-view').innerHTML = selectedOptions.join(', ');
      } else {
        document.querySelector('.select-view').innerHTML = "Выберите жанр";
      }
    }
  });
});

document.addEventListener('click', function (e) {
  if (!document.querySelector('.select-wrapper').contains(e.target)) {
    selectResult.classList.remove('select-active');
  }
});

document.querySelector('.custom-select').addEventListener('change', function (e) {
  document.querySelectorAll('.select-result div').forEach(function (item) {
    item.classList.remove('display-none');
  });
  document.querySelectorAll('.select-result div').forEach(function (item) {
    if (item.innerHTML === e.target.options[e.target.selectedIndex].innerHTML) {
      item.classList.add('display-none');
      document.querySelector('.select-result div').innerHTML = item.innerHTML;
      firstDiv.style.display = 'block';
    }
  });
});