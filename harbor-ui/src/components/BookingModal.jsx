import React, { useState, useEffect } from 'react';

export default function BookingModal({ listing, onClose, onSuccess, showToast }) {
  const [startDay, setStartDay] = useState(null);
  const [endDay, setEndDay] = useState(null);
  const [isBooking, setIsBooking] = useState(false);

  useEffect(() => {
    if (startDay && endDay) {
      if (startDay > endDay) {
        showToast("Дата выезда не может быть раньше заезда", "bg-red-500");
        setEndDay(null);
        return;
      }
      
      let hasOverlap = false;
      for (let d = startDay; d <= endDay; d++) {
        if (listing.bookedDays.includes(d)) {
          hasOverlap = true;
          break;
        }
      }

      if (hasOverlap) {
        showToast("В выбранном диапазоне есть занятые даты!", "bg-red-500");
        setEndDay(null);
      }
    }
  }, [startDay, endDay, listing.bookedDays, showToast]);

  const handleDayClick = (day) => {
    if (listing.bookedDays.includes(day)) return;

    if (!startDay || (startDay && endDay)) {
      setStartDay(day);
      setEndDay(null);
    } else {
      if (day < startDay) {
        setEndDay(startDay);
        setStartDay(day);
      } else {
        setEndDay(day);
      }
    }
  };

  const handleInputChange = (e, type) => {
    const val = e.target.value;
    const day = val ? parseInt(val.split('-')[2], 10) : null;
    
    if (type === 'start') setStartDay(day);
    if (type === 'end') setEndDay(day);
  };

  const executeBooking = async () => {
    if (!listing.id || !startDay || !endDay) return;
    
    try {
      setIsBooking(true);
      showToast("Отправка запроса на Harbor Server...", "bg-blue-500");

      // 1. Генерируем уникальный Correlation ID для цепочки запросов
      const correlationId = "req-" + Date.now();

      // 2. Отправляем запрос с заголовком X-Correlation-ID
      const response = await fetch("http://localhost:8082/book", {
        method: "POST",
        headers: { 
          "Content-Type": "application/json",
          "X-Correlation-ID": correlationId
        },
        body: JSON.stringify({ 
          listing_id: listing.id, 
          user_id: "a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a99"
        }) 
      });

      if (response.ok) {
        showToast("Жилье успешно забронировано!", "bg-green-500");
        onSuccess(listing.id, startDay, endDay);
        setTimeout(onClose, 1500);
      } else {
        showToast("Ошибка сервера при бронировании", "bg-red-500");
      }
    } catch (error) {
      console.error("Ошибка сети:", error);
      showToast("Сервер Harbor недоступен", "bg-red-500");
    } finally {
      setIsBooking(false);
    }
  };

  const renderCalendar = () => {
    const days = [];
    for (let day = 1; day <= 30; day++) {
      const isBooked = listing.bookedDays.includes(day);
      const isSelected = startDay === day || endDay === day;
      const isBetween = startDay && endDay && day > startDay && day < endDay;

      let classes = "p-2 rounded-lg text-center text-xs font-semibold transition-all duration-200 select-none ";

      if (isBooked) {
        classes += "bg-red-50 text-red-400 line-through border border-red-100 cursor-not-allowed opacity-60";
      } else if (isSelected || isBetween) {
        classes += "bg-indigo-600 text-white font-bold shadow-md cursor-pointer transform scale-105";
      } else {
        classes += "bg-green-50 text-green-700 border border-green-100 cursor-pointer hover:bg-green-100";
      }

      days.push(
        <div key={day} className={classes} onClick={() => !isBooked && handleDayClick(day)}>
          {day}
        </div>
      );
    }
    return days;
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50 overflow-y-auto">
      <div className="bg-white rounded-2xl max-w-5xl w-full max-h-[90vh] overflow-hidden shadow-2xl flex flex-col md:flex-row relative">
        
        <button 
          onClick={onClose} 
          className="absolute top-4 right-4 bg-white hover:bg-gray-100 text-gray-700 font-bold p-2 rounded-full shadow-md z-10 w-10 h-10 flex items-center justify-center transition-colors border"
        >
          ✕
        </button>

        <div className="md:w-3/5 p-8 overflow-y-auto border-r border-gray-100">
          <div className="rounded-xl overflow-hidden h-64 mb-6 bg-gray-100">
            <img src={listing.image} alt={listing.title} className="w-full h-full object-cover" />
          </div>
          
          <div className="flex justify-between items-start mb-2">
            <h2 className="text-2xl font-bold text-gray-900">{listing.title}</h2>
            <div className="flex items-center space-x-1 text-yellow-500 font-bold bg-yellow-50 px-2.5 py-1 rounded-lg text-sm border border-yellow-100">
              <span>★</span>
              <span>{listing.rating}</span>
            </div>
          </div>
          
          <p className="text-gray-500 text-sm mb-4">{listing.location} • Макс. гостей: {listing.guests}</p>
          
          <hr className="my-4 border-gray-100" />
          
          <h3 className="text-lg font-bold mb-2 text-gray-800">Описание</h3>
          <p className="text-gray-600 leading-relaxed mb-6 text-sm">{listing.description}</p>
          
          <hr className="my-4 border-gray-100" />
          
          <h3 className="text-lg font-bold mb-3 text-gray-800">Что есть в жилье</h3>
          <div className="grid grid-cols-2 gap-2">
            {listing.amenities.map((amenity, idx) => (
              <div key={idx} className="flex items-center space-x-2 text-sm text-gray-600 bg-gray-50 p-2 rounded-lg">
                <span className="text-indigo-500">✓</span>
                <span>{amenity}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="md:w-2/5 p-8 bg-gray-50 flex flex-col justify-between overflow-y-auto">
          <div>
            <div className="flex justify-between items-baseline mb-6">
              <span className="text-2xl font-black text-gray-900">{listing.price}</span>
              <span className="text-sm text-gray-500">/ ночь</span>
            </div>

            <div className="bg-white p-4 rounded-xl border border-gray-200 shadow-sm mb-6">
              <span className="text-xs font-bold text-gray-400 block mb-3 uppercase tracking-wider">Июнь 2026</span>
              <div className="grid grid-cols-7 gap-1 text-center text-xs font-bold text-gray-400 mb-2">
                <div>Пн</div><div>Вт</div><div>Ср</div><div>Чт</div><div>Пт</div><div>Сб</div><div>Вс</div>
              </div>
              <div className="grid grid-cols-7 gap-1">
                {renderCalendar()}
              </div>
            </div>

            <div className="space-y-4">
              <div className="flex space-x-2">
                <div className="w-1/2">
                  <label className="block text-xs font-bold text-gray-500 mb-1 uppercase">Заезд</label>
                  <input 
                    type="date" 
                    min="2026-06-01" 
                    max="2026-06-30" 
                    value={startDay ? `2026-06-${String(startDay).padStart(2, '0')}` : ''}
                    onChange={(e) => handleInputChange(e, 'start')}
                    className="w-full border border-gray-300 rounded-lg p-2 focus:ring-2 focus:ring-indigo-500 outline-none text-xs bg-white"
                  />
                </div>
                <div className="w-1/2">
                  <label className="block text-xs font-bold text-gray-500 mb-1 uppercase">Выезд</label>
                  <input 
                    type="date" 
                    min="2026-06-01" 
                    max="2026-06-30" 
                    value={endDay ? `2026-06-${String(endDay).padStart(2, '0')}` : ''}
                    onChange={(e) => handleInputChange(e, 'end')}
                    className="w-full border border-gray-300 rounded-lg p-2 focus:ring-2 focus:ring-indigo-500 outline-none text-xs bg-white"
                  />
                </div>
              </div>
              
              <div className="text-center text-xs py-2 min-h-[32px]">
                {!startDay && <span className="text-gray-400">Выберите дату заезда на календаре</span>}
                {startDay && !endDay && <span className="text-indigo-600 font-semibold">Заезд: {startDay} июня. Укажите выезд.</span>}
                {startDay && endDay && <span className="text-green-600 font-bold">Выбран период: {startDay} — {endDay} июня</span>}
              </div>
            </div>
          </div>

          <div className="mt-6">
            <button 
              onClick={executeBooking}
              disabled={!startDay || !endDay || isBooking}
              className="w-full bg-gradient-to-r from-pink-500 to-rose-600 hover:from-pink-600 hover:to-rose-700 text-white font-bold py-3 px-4 rounded-xl transition-all shadow-lg disabled:from-gray-300 disabled:to-gray-300 disabled:cursor-not-allowed disabled:shadow-none transform active:scale-95 text-sm"
            >
              {isBooking ? "Отправка..." : "Забронировать жилье"}
            </button>
            <p className="text-[11px] text-center text-gray-400 mt-2">Пока вы ничего не платите, запрос уйдёт на сервер</p>
          </div>
        </div>

      </div>
    </div>
  );
}