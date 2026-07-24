import { useState } from "react";

interface FileItem {
  id: string;
  name: string;
  type: string; // 'figma' | 'xd' | 'pdf' | 'audio' | 'image' | 'excel' | 'doc'
  date: string;
  size: string;
  category: "documents" | "google_drive" | "one_drive" | "dropbox";
}

export function DashboardPage() {
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Form states for new file
  const [newFileName, setNewFileName] = useState("");
  const [newFileType, setNewFileType] = useState("pdf");
  const [newFileCategory, setNewFileCategory] = useState<"documents" | "google_drive" | "one_drive" | "dropbox">("documents");
  const [newFileSize, setNewFileSize] = useState("");

  const [files, setFiles] = useState<FileItem[]>([
    { id: "1", name: "Xd File", type: "xd", date: "01-03-2021", size: "3.5mb", category: "documents" },
    { id: "2", name: "Figma File", type: "figma", date: "27-02-2021", size: "19mb", category: "google_drive" },
    { id: "3", name: "Documetns", type: "doc", date: "23-02-2021", size: "15mb", category: "documents" },
    { id: "4", name: "Sound File", type: "audio", date: "21-02-2021", size: "40mb", category: "one_drive" },
    { id: "5", name: "Media", type: "image", date: "23-02-2021", size: "15mb", category: "google_drive" },
    { id: "6", name: "Sales PDF", type: "pdf", date: "21-02-2021", size: "9mb", category: "dropbox" },
    { id: "7", name: "Excel File", type: "excel", date: "23-02-2021", size: "11mb", category: "one_drive" },
  ]);

  // Folder configs
  const folderData = {
    documents: { title: "Documents", files: 1328, size: "1.3GB", color: "text-blue-500 bg-blue-50", barColor: "bg-blue-500", progress: 40 },
    google_drive: { title: "Google Drive", files: 2329, size: "2.9GB", color: "text-amber-500 bg-amber-50", barColor: "bg-amber-500", progress: 65 },
    one_drive: { title: "One Drive", files: 1916, size: "1.7GB", color: "text-indigo-500 bg-indigo-50", barColor: "bg-indigo-500", progress: 50 },
    dropbox: { title: "Dropbox", files: 328, size: "1.1GB", color: "text-cyan-500 bg-cyan-50", barColor: "bg-cyan-500", progress: 25 },
  };

  // Add new file handler
  const handleAddFile = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newFileName || !newFileSize) return;

    const newFile: FileItem = {
      id: Date.now().toString(),
      name: newFileName,
      type: newFileType,
      date: new Date().toLocaleDateString("en-GB").replace(/\//g, "-"),
      size: newFileSize.toLowerCase().endsWith("mb") ? newFileSize.toLowerCase() : `${newFileSize}mb`,
      category: newFileCategory,
    };

    setFiles([newFile, ...files]);
    setIsModalOpen(false);

    // Reset inputs
    setNewFileName("");
    setNewFileSize("");
  };

  // Filter files
  const filteredFiles = files.filter((f) => {
    const matchesCategory = selectedCategory ? f.category === selectedCategory : true;
    const matchesSearch = f.name.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesCategory && matchesSearch;
  });

  const getFileIcon = (type: string) => {
    switch (type) {
      case "xd":
        return (
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-pink-50 text-pink-500">
            <span className="font-bold text-xs">Xd</span>
          </div>
        );
      case "figma":
        return (
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-purple-50 text-purple-500">
            <span className="font-bold text-xs">Fg</span>
          </div>
        );
      case "doc":
      case "pdf":
        return (
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-red-50 text-red-500">
            <span className="font-bold text-xs">PDF</span>
          </div>
        );
      case "audio":
        return (
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-orange-50 text-orange-500">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
            </svg>
          </div>
        );
      case "image":
        return (
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 text-blue-500">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </div>
        );
      case "excel":
        return (
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-green-50 text-green-500">
            <span className="font-bold text-xs">XL</span>
          </div>
        );
      default:
        return (
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-50 text-slate-500">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
          </div>
        );
    }
  };

  return (
    <div className="flex flex-col gap-8 animate-fadeIn">
      {/* Search and Action Bar */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-bold font-display text-slate-800">My Files</h2>
          <p className="text-xs font-semibold text-slate-400 mt-1">
            {selectedCategory ? `Viewing ${folderData[selectedCategory as keyof typeof folderData].title}` : "Viewing all cloud storage folders"}
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          {/* Quick search input */}
          <div className="relative flex items-center bg-white border border-slate-200/60 rounded-xl px-3 py-1.5 shadow-sm">
            <svg className="w-4 h-4 text-slate-400 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              type="text"
              placeholder="Search files…"
              className="text-xs bg-transparent outline-none w-36 sm:w-44 text-slate-700 placeholder:text-slate-400"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            {searchQuery && (
              <button onClick={() => setSearchQuery("")} className="text-slate-400 hover:text-slate-600 ml-1">
                <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            )}
          </div>

          {selectedCategory && (
            <button
              onClick={() => setSelectedCategory(null)}
              className="btn bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-xl px-3 py-2 text-xs"
            >
              Clear Filter
            </button>
          )}

          <button
            onClick={() => setIsModalOpen(true)}
            className="btn-primary rounded-xl px-4 py-2.5 text-xs inline-flex items-center gap-1.5 font-bold shadow-lg shadow-blue-500/10 hover:shadow-blue-500/25 active:scale-95 transition-all"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
            </svg>
            Add New
          </button>
        </div>
      </div>

      {/* Grid of Folders */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        {Object.entries(folderData).map(([key, f]) => {
          const isActive = selectedCategory === key;
          return (
            <div
              key={key}
              onClick={() => setSelectedCategory(isActive ? null : key)}
              className={`card flex flex-col justify-between cursor-pointer border p-5 transition-all duration-300 ${
                isActive
                  ? "border-blue-500 bg-blue-50/20 shadow-md ring-2 ring-blue-500/20"
                  : "border-slate-100/80 hover:border-slate-300/80 hover:scale-[1.02]"
              }`}
            >
              <div className="flex items-center justify-between">
                <div className={`flex h-11 w-11 items-center justify-center rounded-xl ${f.color}`}>
                  <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                  </svg>
                </div>
                <button className="text-slate-300 hover:text-slate-500">
                  <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                  </svg>
                </button>
              </div>

              <div className="mt-5">
                <h3 className="font-display font-bold text-slate-800 text-[15px]">{f.title}</h3>
                
                {/* Progress bar */}
                <div className="mt-3.5 h-1.5 w-full rounded-full bg-slate-100 overflow-hidden">
                  <div className={`h-full ${f.barColor} rounded-full`} style={{ width: `${f.progress}%` }} />
                </div>

                <div className="mt-3 flex items-center justify-between text-[11px] font-bold text-slate-400">
                  <span>{f.files.toLocaleString()} Files</span>
                  <span className="text-slate-600">{f.size}</span>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Columns: Left Table & Right storage details */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Left Column (Recent Files & Analytics) */}
        <div className="lg:col-span-2 flex flex-col gap-6">
          {/* Recent Files */}
          <div className="card border border-slate-100/80 p-6 flex flex-col gap-4">
            <div className="flex items-center justify-between border-b border-slate-100 pb-3">
              <h3 className="font-display font-bold text-slate-800 text-[16px]">Recent Files</h3>
              <button className="text-xs font-bold text-slate-400 hover:text-blue-500 transition-colors">See more</button>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100/60 pb-2">
                    <th className="py-2.5 font-bold">File Name</th>
                    <th className="py-2.5 font-bold">Date</th>
                    <th className="py-2.5 font-bold">Size</th>
                    <th className="py-2.5"></th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-50">
                  {filteredFiles.length === 0 ? (
                    <tr>
                      <td colSpan={4} className="py-8 text-center text-sm text-slate-400">
                        No files found matching the search.
                      </td>
                    </tr>
                  ) : (
                    filteredFiles.map((file) => (
                      <tr key={file.id} className="group hover:bg-slate-50/50 transition-colors">
                        <td className="py-3 flex items-center gap-3">
                          {getFileIcon(file.type)}
                          <div className="flex flex-col">
                            <span className="text-sm font-bold text-slate-700 group-hover:text-blue-500 transition-colors leading-snug">{file.name}</span>
                            <span className="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">{file.type} file</span>
                          </div>
                        </td>
                        <td className="py-3 text-xs font-semibold text-slate-500">{file.date}</td>
                        <td className="py-3 text-xs font-bold text-slate-600">{file.size}</td>
                        <td className="py-3 text-right">
                          <button
                            onClick={() => setFiles(files.filter((f) => f.id !== file.id))}
                            className="text-slate-300 hover:text-red-500 p-1 opacity-0 group-hover:opacity-100 transition-all rounded-lg hover:bg-red-50"
                            title="Delete file"
                          >
                            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                              <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                          </button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>

          {/* Analytics Chart */}
          <div className="card border border-slate-100/80 p-6 flex flex-col gap-4">
            <div className="flex items-center justify-between border-b border-slate-100 pb-3">
              <h3 className="font-display font-bold text-slate-800 text-[16px]">Analytics</h3>
              <button className="text-slate-300 hover:text-slate-500">
                <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z" />
                </svg>
              </button>
            </div>

            {/* Visual HTML/CSS Chart */}
            <div className="flex gap-4 items-end mt-4 h-56">
              {/* Y Axis scale */}
              <div className="flex flex-col justify-between h-48 text-[10px] font-bold text-slate-400 w-6">
                <span>100</span>
                <span>80</span>
                <span>60</span>
                <span>40</span>
                <span>20</span>
                <span>0</span>
              </div>

              {/* Columns */}
              <div className="flex-1 flex justify-around items-end h-48 border-b border-slate-100 pb-1">
                {/* Sat */}
                <div className="flex flex-col items-center gap-2 group cursor-pointer">
                  <div className="flex gap-1 items-end h-32 relative">
                    <div className="w-3 bg-blue-400 rounded-full h-[65%] hover:opacity-85 transition-all" title="Uploads: 65%" />
                    <div className="w-3 bg-blue-600 rounded-full h-[85%] hover:opacity-85 transition-all" title="Downloads: 85%" />
                  </div>
                  <span className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider">Sat</span>
                </div>

                {/* Sun */}
                <div className="flex flex-col items-center gap-2 group cursor-pointer">
                  <div className="flex gap-1 items-end h-32">
                    <div className="w-3 bg-blue-400 rounded-full h-[40%] hover:opacity-85 transition-all" title="Uploads: 40%" />
                    <div className="w-3 bg-blue-600 rounded-full h-[60%] hover:opacity-85 transition-all" title="Downloads: 60%" />
                  </div>
                  <span className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider">Sun</span>
                </div>

                {/* Mon */}
                <div className="flex flex-col items-center gap-2 group cursor-pointer">
                  <div className="flex gap-1 items-end h-32">
                    <div className="w-3 bg-blue-400 rounded-full h-[90%] hover:opacity-85 transition-all" title="Uploads: 90%" />
                    <div className="w-3 bg-blue-600 rounded-full h-[75%] hover:opacity-85 transition-all" title="Downloads: 75%" />
                  </div>
                  <span className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider">Mon</span>
                </div>

                {/* Tue */}
                <div className="flex flex-col items-center gap-2 group cursor-pointer">
                  <div className="flex gap-1 items-end h-32">
                    <div className="w-3 bg-blue-400 rounded-full h-[45%] hover:opacity-85 transition-all" title="Uploads: 45%" />
                    <div className="w-3 bg-blue-600 rounded-full h-[35%] hover:opacity-85 transition-all" title="Downloads: 35%" />
                  </div>
                  <span className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider">Tue</span>
                </div>

                {/* Wed */}
                <div className="flex flex-col items-center gap-2 group cursor-pointer">
                  <div className="flex gap-1 items-end h-32">
                    <div className="w-3 bg-blue-400 rounded-full h-[95%] hover:opacity-85 transition-all" title="Uploads: 95%" />
                    <div className="w-3 bg-blue-600 rounded-full h-[85%] hover:opacity-85 transition-all" title="Downloads: 85%" />
                  </div>
                  <span className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider">Wed</span>
                </div>

                {/* Thu */}
                <div className="flex flex-col items-center gap-2 group cursor-pointer">
                  <div className="flex gap-1 items-end h-32">
                    <div className="w-3 bg-amber-400 rounded-full h-[75%] hover:opacity-85 transition-all" title="Uploads: 75%" />
                    <div className="w-3 bg-amber-500 rounded-full h-[65%] hover:opacity-85 transition-all" title="Downloads: 65%" />
                  </div>
                  <span className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider">Thu</span>
                </div>

                {/* Fri */}
                <div className="flex flex-col items-center gap-2 group cursor-pointer">
                  <div className="flex gap-1 items-end h-32">
                    <div className="w-3 bg-blue-400 rounded-full h-[55%] hover:opacity-85 transition-all" title="Uploads: 55%" />
                    <div className="w-3 bg-blue-600 rounded-full h-[45%] hover:opacity-85 transition-all" title="Downloads: 45%" />
                  </div>
                  <span className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider">Fri</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column (Storage Details) */}
        <div className="flex flex-col gap-6">
          <div className="card border border-slate-100/80 p-6 flex flex-col gap-6 h-full justify-between">
            <div className="flex items-center justify-between border-b border-slate-100 pb-3">
              <h3 className="font-display font-bold text-slate-800 text-[16px]">Storage Details</h3>
              <button className="text-slate-300 hover:text-slate-500">
                <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                </svg>
              </button>
            </div>

            {/* Doughnut Chart */}
            <div className="flex justify-center py-4">
              <div className="relative w-44 h-44 flex items-center justify-center">
                {/* SVG Segmented Circle */}
                <svg className="w-full h-full transform -rotate-90" viewBox="0 0 120 120">
                  {/* Background Track */}
                  <circle cx="60" cy="60" r="46" fill="transparent" stroke="#f1f5f9" strokeWidth="12" />
                  
                  {/* Segment 1: Media Files (Cyan) - 52% */}
                  <circle
                    cx="60"
                    cy="60"
                    r="46"
                    fill="transparent"
                    stroke="#22d3ee"
                    strokeWidth="14"
                    strokeDasharray="289"
                    strokeDashoffset="138" // (1 - 0.52) * 289
                    strokeLinecap="round"
                  />
                  {/* Segment 2: Other Files (Yellow) - 44% */}
                  <circle
                    cx="60"
                    cy="60"
                    r="46"
                    fill="transparent"
                    stroke="#fbbf24"
                    strokeWidth="14"
                    strokeDasharray="289"
                    strokeDashoffset="162" // Rotated/drawn relative
                    className="transform rotate-[187deg] origin-center"
                    strokeLinecap="round"
                  />
                  {/* Segment 3: Documents Files (Blue) - 4.5% */}
                  <circle
                    cx="60"
                    cy="60"
                    r="46"
                    fill="transparent"
                    stroke="#3b82f6"
                    strokeWidth="14"
                    strokeDasharray="289"
                    strokeDashoffset="276"
                    className="transform rotate-[345deg] origin-center"
                    strokeLinecap="round"
                  />
                  {/* Segment 4: Unknown Files (Red) - 4.5% */}
                  <circle
                    cx="60"
                    cy="60"
                    r="46"
                    fill="transparent"
                    stroke="#ef4444"
                    strokeWidth="14"
                    strokeDasharray="289"
                    strokeDashoffset="276"
                    className="transform rotate-[361deg] origin-center"
                    strokeLinecap="round"
                  />
                </svg>

                {/* Inner Text */}
                <div className="absolute inset-0 flex flex-col items-center justify-center">
                  <span className="text-3xl font-display font-extrabold text-slate-800 leading-none">29.1</span>
                  <span className="text-[10px] font-bold text-slate-400 mt-1 uppercase tracking-wider">Of 128GB</span>
                </div>
              </div>
            </div>

            {/* Legend list */}
            <div className="flex flex-col gap-3">
              {/* Documents Files */}
              <div className="flex items-center justify-between border border-slate-100 rounded-2xl p-3 hover:bg-slate-50 transition-colors">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-blue-50 text-blue-500">
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <div className="flex flex-col">
                    <span className="text-xs font-bold text-slate-700">Documents Files</span>
                    <span className="text-[9px] font-bold text-slate-400 mt-0.5">1,328 Files</span>
                  </div>
                </div>
                <span className="text-xs font-bold text-slate-600">1.3GB</span>
              </div>

              {/* Media Files */}
              <div className="flex items-center justify-between border border-slate-100 rounded-2xl p-3 hover:bg-slate-50 transition-colors">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-cyan-50 text-cyan-500">
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-.553.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <div className="flex flex-col">
                    <span className="text-xs font-bold text-slate-700">Media Files</span>
                    <span className="text-[9px] font-bold text-slate-400 mt-0.5">1,328 Files</span>
                  </div>
                </div>
                <span className="text-xs font-bold text-slate-600">15.1GB</span>
              </div>

              {/* Other Files */}
              <div className="flex items-center justify-between border border-slate-100 rounded-2xl p-3 hover:bg-slate-50 transition-colors">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-amber-50 text-amber-500">
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
                    </svg>
                  </div>
                  <div className="flex flex-col">
                    <span className="text-xs font-bold text-slate-700">Other Files</span>
                    <span className="text-[9px] font-bold text-slate-400 mt-0.5">1,328 Files</span>
                  </div>
                </div>
                <span className="text-xs font-bold text-slate-600">12.7GB</span>
              </div>

              {/* Unknown Files */}
              <div className="flex items-center justify-between border border-slate-100 rounded-2xl p-3 hover:bg-slate-50 transition-colors">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-red-50 text-red-500">
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <div className="flex flex-col">
                    <span className="text-xs font-bold text-slate-700">Unknown</span>
                    <span className="text-[9px] font-bold text-slate-400 mt-0.5">428 Files</span>
                  </div>
                </div>
                <span className="text-xs font-bold text-slate-600">1.3GB</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Add New File Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm animate-fadeIn">
          <div className="bg-white rounded-3xl border border-slate-100 shadow-2xl p-6 w-full max-w-md flex flex-col gap-4 transform transition-all scale-100">
            <div className="flex items-center justify-between border-b border-slate-100 pb-3">
              <h3 className="font-display font-bold text-slate-800 text-[18px]">Upload New File</h3>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-slate-400 hover:text-slate-600 p-1 rounded-lg hover:bg-slate-50 transition-colors"
              >
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <form onSubmit={handleAddFile} className="flex flex-col gap-4">
              <div>
                <label className="label">File Name</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Project Proposal"
                  className="input"
                  value={newFileName}
                  onChange={(e) => setNewFileName(e.target.value)}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">File Type</label>
                  <select
                    className="input"
                    value={newFileType}
                    onChange={(e) => setNewFileType(e.target.value)}
                  >
                    <option value="pdf">PDF File</option>
                    <option value="figma">Figma Design</option>
                    <option value="xd">Adobe XD</option>
                    <option value="doc">Word Document</option>
                    <option value="audio">Sound Track</option>
                    <option value="image">Image Media</option>
                    <option value="excel">Excel Sheet</option>
                  </select>
                </div>

                <div>
                  <label className="label">File Size</label>
                  <input
                    type="text"
                    required
                    placeholder="e.g. 15mb"
                    className="input"
                    value={newFileSize}
                    onChange={(e) => setNewFileSize(e.target.value)}
                  />
                </div>
              </div>

              <div>
                <label className="label">Cloud Location</label>
                <select
                  className="input"
                  value={newFileCategory}
                  onChange={(e) => setNewFileCategory(e.target.value as any)}
                >
                  <option value="documents">Local Documents Folder</option>
                  <option value="google_drive">Google Drive</option>
                  <option value="one_drive">One Drive</option>
                  <option value="dropbox">Dropbox</option>
                </select>
              </div>

              <div className="flex gap-3 justify-end mt-2 pt-2 border-t border-slate-100">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="btn bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-xl py-2 px-4 text-xs font-bold"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="btn bg-[#2F80ED] hover:bg-blue-600 text-white rounded-xl py-2 px-4 text-xs font-bold shadow-lg shadow-blue-500/10"
                >
                  Add File
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
