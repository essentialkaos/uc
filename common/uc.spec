################################################################################

%define debug_package  %{nil}

################################################################################

Summary:        Simple utility for counting unique lines
Name:           uc
Version:        3.1.1
Release:        0%{?dist}
Group:          Applications/System
License:        Apache License, Version 2.0
URL:            https://kaos.sh/uc

Source0:        https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  golang >= 1.24

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
Simple utility for counting unique lines.

################################################################################

%prep

%setup -q
if [[ ! -d "%{name}/vendor" ]] ; then
  echo -e "----\nThis package requires vendored dependencies\n----"
  exit 1
elif [[ -f "%{name}/%{name}" ]] ; then
  echo -e "----\nSources must not contain precompiled binaries\n----"
  exit 1
fi

%build
pushd %{name}
  go build %{name}.go
  cp LICENSE ..
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -dm 755 %{buildroot}%{_mandir}/man1

install -pm 755 %{name}/%{name} %{buildroot}%{_bindir}/

./%{name}/%{name} --generate-man > %{buildroot}%{_mandir}/man1/%{name}.1

%clean
rm -rf %{buildroot}

%post
if [[ -d %{_sysconfdir}/bash_completion.d ]] ; then
  %{name} --completion=bash 1> %{_sysconfdir}/bash_completion.d/%{name} 2>/dev/null
fi

if [[ -d %{_datarootdir}/fish/vendor_completions.d ]] ; then
  %{name} --completion=fish 1> %{_datarootdir}/fish/vendor_completions.d/%{name}.fish 2>/dev/null
fi

if [[ -d %{_datadir}/zsh/site-functions ]] ; then
  %{name} --completion=zsh 1> %{_datadir}/zsh/site-functions/_%{name} 2>/dev/null
fi

%postun
if [[ $1 == 0 ]] ; then
  if [[ -f %{_sysconfdir}/bash_completion.d/%{name} ]] ; then
    rm -f %{_sysconfdir}/bash_completion.d/%{name} &>/dev/null || :
  fi

  if [[ -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish ]] ; then
    rm -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish &>/dev/null || :
  fi

  if [[ -f %{_datadir}/zsh/site-functions/_%{name} ]] ; then
    rm -f %{_datadir}/zsh/site-functions/_%{name} &>/dev/null || :
  fi
fi

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_mandir}/man1/%{name}.1.*
%{_bindir}/%{name}

################################################################################

%changelog
* Wed Oct 15 2025 Anton Novojilov <andy@essentialkaos.com> - 3.1.1-0
- Fixed bug with handling self-update option

* Wed Oct 15 2025 Anton Novojilov <andy@essentialkaos.com> - 3.1.0-0
- Code refactoring
- Dependencies update

* Fri May 16 2025 Anton Novojilov <andy@essentialkaos.com> - 3.0.3-0
- Code refactoring
- Dependencies update

* Mon Jun 24 2024 Anton Novojilov <andy@essentialkaos.com> - 3.0.2-0
- Code refactoring
- Dependencies update

* Thu Mar 28 2024 Anton Novojilov <andy@essentialkaos.com> - 3.0.1-0
- Improved support information gathering
- Code refactoring
- Dependencies update

* Mon Feb 19 2024 Anton Novojilov <andy@essentialkaos.com> - 3.0.0-0
- crc64 replaced by xxhash
- Added different output formats for distribution info
- Code refactoring
- Dependencies update

* Tue Dec 19 2023 Anton Novojilov <andy@essentialkaos.com> - 2.0.1-0
- Dependencies update
- Code refactoring

* Wed May 03 2023 Anton Novojilov <andy@essentialkaos.com> - 2.0.0-0
- Better standard input handling

* Mon Mar 06 2023 Anton Novojilov <andy@essentialkaos.com> - 1.1.1-0
- Added verbose info output
- Dependencies update
- Code refactoring

* Thu Dec 01 2022 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-1
- Fixed build using sources from source.kaos.st

* Wed Aug 10 2022 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- Minor UI improvements
- Fixed bug with parsing max number of unique lines
- Updated dependencies

* Tue Mar 29 2022 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Removed pkg.re usage
- Added module info
- Added Dependabot configuration

* Thu Oct 22 2020 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Added possibility to define -m/--max option as number with K and M
- Added man page generation

* Fri Jan 31 2020 Anton Novojilov <andy@essentialkaos.com> - 0.0.2-0
- Added option -m/--max for defining maximum unique lines to process
- Added option -d/--dist for printing data distribution

* Thu Jan 30 2020 Anton Novojilov <andy@essentialkaos.com> - 0.0.1-0
- Initial build
