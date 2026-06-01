import React from 'react';

export default function ListingCard({ listing, onClick }) {
  return (
    <div 
      className="bg-white rounded-2xl shadow-sm overflow-hidden hover:shadow-md transition-shadow duration-300 cursor-pointer border border-gray-100"
      onClick={() => onClick(listing)}
    >
      {/* Контейнер для картинки с ограничением overflow, чтобы при увеличении края не вылезали */}
      <div className="relative h-56 w-full overflow-hidden bg-gray-100">
        <img 
          src={listing.image} 
          alt={listing.title} 
          className="w-full h-full object-cover transition-transform duration-500 ease-out hover:scale-107"
        />
      </div>
      <div className="p-6">
        <div className="flex justify-between items-start mb-1">
          <h2 className="text-base font-bold text-gray-800 line-clamp-1">{listing.title}</h2>
          <div className="flex items-center space-x-0.5 text-sm font-bold text-gray-700 shrink-0">
            <span>★</span>
            <span>{listing.rating}</span>
          </div>
        </div>
        <p className="text-gray-400 text-xs mb-3">{listing.location}</p>
        <p className="text-gray-900 text-sm font-black">{listing.price} <span className="text-gray-400 font-normal text-xs">/ ночь</span></p>
      </div>
    </div>
  );
}