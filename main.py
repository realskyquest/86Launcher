from email import message
import tkinter as tk
from tkinter import ttk, filedialog, messagebox
from turtle import down
from PIL import Image, ImageTk
import requests
import threading
import requests
import os

class Engine:
    def __init__(self, root):
        print("Class Engine, init is executed")

        self.root = root
        self.root.title("86 Game Launcher")
        self.root.geometry("450x520")

        self.is_downloading = False

        self.CreateUI()
        print("UI has been drawn")
    
    def CreateUI(self):
        self.image_frame = tk.Frame(self.root)
        self.image_frame.pack(pady=10)
        
        self.image = Image.open("banner.jpg")
        self.image = self.image.resize((300, 100))
        self.photo = ImageTk.PhotoImage(self.image)
        self.image_label = tk.Label(self.image_frame, image=self.photo) # type: ignore
        self.image_label.pack()
        print("Banner loaded")
    
        # Combobox
        self.version_combobox = tk.Frame(self.root)
        self.version_combobox.pack(pady=8)

        self.version_combobox_msg_label = tk.Label(self.version_combobox, text="Select game version:")
        self.version_combobox_msg_label.pack(side=tk.LEFT)

        self.version_combobox_var = tk.StringVar()
        self.version_combobox_menu = ttk.Combobox(self.version_combobox, textvariable=self.version_combobox_var, state="readonly")
        self.version_combobox_menu.pack(side=tk.LEFT)

        self.GithubReleasesThread()

        # Download location
        self.download_location_frame = tk.Frame(self.root)
        self.download_location_frame.pack(pady=4)
        
        self.download_location_label = tk.Label(self.download_location_frame, text="Download path:")
        self.download_location_label.pack(side=tk.LEFT)
        
        self.download_location_path_var = tk.StringVar()
        self.download_location_entry = tk.Entry(self.download_location_frame, textvariable=self.download_location_path_var, width=30)
        self.download_location_entry.pack(side=tk.LEFT)

        self.browse_location_button = tk.Button(self.download_location_frame, text="Browse", command=self.BrowseDownloadLocation)
        self.browse_location_button.pack(side=tk.LEFT)
        
        # Get button
        self.get_button = tk.Button(self.root, text="Get game assets", command=self.GetGame)
        self.get_button.pack(pady=8)

        # Select asset
        self.assets_combobox = tk.Frame(self.root)
        self.assets_combobox.pack(pady=8)

        self.assets_combobox_msg_label = tk.Label(self.assets_combobox, text="Select asset:")
        self.assets_combobox_msg_label.pack(side=tk.LEFT)

        self.assets_combobox_var = tk.StringVar()
        self.assets_combobox_menu = ttk.Combobox(self.assets_combobox, textvariable=self.assets_combobox_var, state="readonly")
        self.assets_combobox_menu.pack(side=tk.LEFT)

        # Refresh assets
        self.asset_select_button = tk.Button(self.root, text="Check asset", command=self.GetAsset)
        self.asset_select_button.pack(pady=8)

        # About asset
        self.asset_info = tk.Frame(self.root)
        self.asset_info.pack()

        self.asset_info_name_label = tk.Label(self.asset_info, text="Name = ")
        self.asset_info_name_label.pack()

        self.asset_info_size_label = tk.Label(self.asset_info, text="Total size = ")
        self.asset_info_size_label.pack()

        self.asset_info_url_label = tk.Label(self.asset_info, text="Link = ")
        self.asset_info_url_label.pack()

        # Download asset
        self.download_button = tk.Button(self.root, text="Download asset", command=self.DownloadAsset)
        self.download_button.pack(pady=8)

        # Download stats
        self.download_progress_var = tk.IntVar()
        self.download_progress = ttk.Progressbar(self.root, orient="horizontal", length=300, mode="determinate", variable=self.download_progress_var, maximum=300)
        self.download_progress.pack(pady=8)

        self.download_progress_label = tk.Label(self.root, text="Downloaded: 0 MiB / Total 0 MiB")
        self.download_progress_label.pack(pady=2)

    def GetAsset(self):
        selected_asset = self.assets_combobox_var.get()
        for i in range(len(self.assets)):
            if selected_asset == self.assets[i]["name"]:
                self.asset_file_size = self.assets[i]["size"]
                self.asset_file_url = self.assets[i]["browser_download_url"]
                self.asset_file_name = self.assets[i]["name"]

                print(self.asset_file_size, self.asset_file_url, self.asset_file_name)

                self.asset_info_name_label.config(text=f"Total size = {self.asset_file_name}")
                self.asset_info_size_label.config(text=f"Total size = {self.asset_file_size / 1048576: .2f} MiB")
                self.asset_info_url_label.config(text=f"Total size = {self.asset_file_url}")

    def BrowseDownloadLocation(self):
        download_path = filedialog.askdirectory()
        if download_path:
            print("Download path is valid")
            self.download_location_path_var.set(download_path)
        else:
            print("Download path is invalid")

    def GithubReleasesThread(self):
        repo_url = "https://api.github.com/repos/Taliayaya/Project-86/releases"
        response = requests.get(repo_url)
        
        if response.status_code == 200:
            releases = response.json()
            versions = [release['tag_name'] for release in releases]
            self.version_combobox_menu['values'] = versions
            if versions:
                self.version_combobox_menu.current(0)
                print("Successfully obtained versions")
        else:
            print("Failed to get versions, please restart the application")
            messagebox.showerror("Error", "Failed to fetch releases from GitHub")
    
    def GetGame(self):
        print("Get button has been pressed")
        getThread = threading.Thread(target=self.GetGameThread)   
        getThread.start() 

    def DownloadAsset(self):
        if self.is_downloading == False:
            self.is_downloading = True
            print("Download button has been pressed")
            downloadThread = threading.Thread(target=self.DownloadAssetThread)
            downloadThread.start()
        else:
            messagebox.showerror("ERROR", "Download is in progress, please quit app to cancel")

    def GetGameThread(self):
        selected_version = self.version_combobox_var.get()
        download_path = self.download_location_path_var.get()

        if not selected_version:
            print("Game version was not selected")
            messagebox.showerror("Error", "Please select a game version")
            return
        
        if not download_path:
            print("Download path was not selected")
            messagebox.showerror("Error", "Please select a download path")
            return

        repo_url = f"https://api.github.com/repos/Taliayaya/Project-86/releases/tags/{selected_version}"
        response = requests.get(repo_url)
        
        if response.status_code == 200:
            release = response.json()
            self.assets = release.get('assets', [])
            
            if not self.assets:
                print("No assets found")
                messagebox.showerror("Error", "No assets found for the selected release")
                return
            
            assets_list = [asset["name"] for asset in self.assets]
            self.assets_combobox_menu["values"] = assets_list
            if assets_list:
                self.assets_combobox_menu.current(0)
            
        else:
            print("Failed to get selected release")
            messagebox.showerror("Error", "Failed to fetch the selected release from GitHub")
    
    def DownloadAssetThread(self):
        download_path = self.download_location_path_var.get()

        response = requests.get(self.asset_file_url, stream=True)
        self.bytes_downloaded = 0

        if response.status_code == 200:
            with open(f"{download_path}/{self.asset_file_name}", "wb") as f:
                for chunk in response.iter_content(chunk_size=8192):
                    f.write(chunk)

                    self.bytes_downloaded += len(chunk)
                    self.download_progress_label.config(text=f"Downloaded: {self.bytes_downloaded / 1048576: .2f} MiB / Total: {self.asset_file_size / 1048576: .2f} MiB")
                    print(f"{self.bytes_downloaded / 1048576: .2f} MiB / Total: {self.asset_file_size / 1048576: .2f} MiB")

                    if self.asset_file_size != 0:
                        ratio = self.bytes_downloaded / self.asset_file_size
                        scaled_value = int(ratio * 300)
                        print(scaled_value)
                        self.download_progress_var.set(scaled_value)

            print("File downloaded successfully!")
            messagebox.showinfo("SUCCESS", "Asset has been downloaded")
            self.is_downloading = False
        else:
            print("Failed to download file.")

def main():
    print("Main Function has been executed")
    root = tk.Tk()
    Engine(root)
    root.mainloop()

if __name__ == '__main__':
    main()

