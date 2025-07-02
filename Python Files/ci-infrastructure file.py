cat ci-infrastructure.yml
---
- name: deploy Jenkins
  hosts: worker
  become: yes

  tasks:

   # Show the inventory path
    - name: show env
      ansible.builtin.debug:
        var: ansible_inventory_sources

 # Show the inventory limit variable
    - name: show ansible_limit
      ansible.builtin.debug:
        var: ansible_limit

 # Show the play hosts
    - name: show ansible hosts
      ansible.builtin.debug:
        var: ansible_play_hosts_all

# Show the environment
    - name: Show env
      ansible.builtin.debug:
        var: ansible_env

    - name: Ensure curl is installed
      ansible.builtin.apt:
        name: curl
        state: present
        update_cache: yes

    - name: Add Jenkins GPG key to trusted keyring (using curl)
      ansible.builtin.shell: "curl -fsSL https://pkg.jenkins.io/debian/jenkins.io-2023.key | tee /usr/share/keyrings/jenkins-keyring.asc > /dev/null"

    - name: Create Jenkins repository file
      ansible.builtin.copy:
        dest: /etc/apt/sources.list.d/jenkins.list
        content: "deb [arch=amd64 signed-by=/usr/share/keyrings/jenkins-keyring.asc] http://pkg.jenkins.io/debian-stable binary/\n"
        mode: '0644'

    - name: Add Jenkins repository GPG key for apt
      ansible.builtin.apt_key:
        file: /usr/share/keyrings/jenkins-keyring.asc
        state: present

    - name: Update apt cache
      ansible.builtin.apt:
        update_cache: yes

    - name: Install Java
      ansible.builtin.apt:
        name: 'openjdk-17-jre'
        state: present
        update_cache: true

    - name: Start Jenkins service
      ansible.builtin.systemd:
        name: jenkins
        state: started
        enabled: true
      register: start_jenkins

    - name: Show start
      ansible.builtin.debug:
        var: start_jenkins

    - name: Get Jenkins status
      ansible.builtin.systemd:
        name: jenkins
        state: started
      register: jenkins_status
      failed_when: jenkins_status.status.ActiveState != 'active'

    - name: Show Jenkins status
      ansible.builtin.debug:
        var: jenkins_status