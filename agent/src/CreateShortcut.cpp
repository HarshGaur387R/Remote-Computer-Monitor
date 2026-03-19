#include <windows.h>
#include <shobjidl.h>
#include <shlguid.h>
#include <objbase.h>
#include <propvarutil.h>
#include <propkey.h>
#include <string>

bool CreateShortcut(const std::wstring& shortcutPath,
                    const std::wstring& exePath,
                    const std::wstring& appId) {
    HRESULT hr = CoInitializeEx(NULL, COINIT_APARTMENTTHREADED);
    if (FAILED(hr)) return false;

    IShellLinkW* pShellLink = nullptr;  // <-- use W version
    hr = CoCreateInstance(CLSID_ShellLink, NULL, CLSCTX_INPROC_SERVER,
                          IID_IShellLinkW, (LPVOID*)&pShellLink);
    if (FAILED(hr)) {
        CoUninitialize();
        return false;
    }

    // Now SetPath accepts LPCWSTR
    pShellLink->SetPath(exePath.c_str());

    IPropertyStore* pPropStore = nullptr;
    hr = pShellLink->QueryInterface(IID_IPropertyStore, (LPVOID*)&pPropStore);
    if (SUCCEEDED(hr)) {
        PROPVARIANT pv;
        hr = InitPropVariantFromString(appId.c_str(), &pv);
        if (SUCCEEDED(hr)) {
            pPropStore->SetValue(PKEY_AppUserModel_ID, pv);
            pPropStore->Commit();
            PropVariantClear(&pv);
        }
        pPropStore->Release();
    }

    IPersistFile* pPersistFile = nullptr;
    hr = pShellLink->QueryInterface(IID_IPersistFile, (LPVOID*)&pPersistFile);
    if (SUCCEEDED(hr)) {
        hr = pPersistFile->Save(shortcutPath.c_str(), TRUE);
        pPersistFile->Release();
    }

    pShellLink->Release();
    CoUninitialize();
    return SUCCEEDED(hr);
}
